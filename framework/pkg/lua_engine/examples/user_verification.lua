-- 用户实名验证规则脚本
-- 基于UserDTO结构进行用户实名状态验证

-- 获取用户信息
local user = context.user

-- 基础验证：用户是否存在
if not user then
    valid = false
    action = "reject"
    error = "用户信息不存在"
    error_reason = "user.not_found"
    variables = {
        verification_result = {
            status = "failed",
            message = "用户信息不存在"
        }
    }
    return
end

-- 1. 实名认证状态验证
local isRealNameVerified = user.isReal == 1
if not isRealNameVerified then
    valid = false
    action = "require_real_name"
    error = "用户未完成实名认证"
    error_reason = "user.real_name_required"
    variables = {
        verification_result = {
            status = "failed",
            message = "用户未完成实名认证",
            required_action = "real_name_verification",
            user_info = {
                id = user.id,
                nickname = user.nickname,
                phone = user.phone,
                account = user.account,
                is_real = user.isReal,
                real_at = user.realAt
            }
        }
    }
    return
end

-- 2. 实名认证时间验证（实名认证必须在7天内完成）
local currentTime = now()
local realNameTime = user.realAt
local timeDiff = currentTime - realNameTime
local daysDiff = timeDiff / (24 * 3600) -- 转换为天数

if daysDiff > 7 then
    valid = false
    action = "reject"
    error = "实名认证时间超过7天限制"
    error_reason = "user.real_name_expired"
    variables = {
        verification_result = {
            status = "failed",
            message = "实名认证时间超过7天限制",
            time_info = {
                real_name_time = realNameTime,
                current_time = currentTime,
                days_diff = daysDiff
            }
        }
    }
    return
end

-- 3. 用户状态验证
local userStatus = user.status
if userStatus ~= 1 then
    valid = false
    action = "reject"
    error = "用户状态异常"
    error_reason = "user.status_invalid"
    variables = {
        verification_result = {
            status = "failed",
            message = "用户状态异常",
            user_status = userStatus
        }
    }
    return
end

-- 4. 年龄验证（必须年满18岁）
local userAge = user.age
if userAge < 18 then
    valid = false
    action = "reject"
    error = "用户年龄未满18岁"
    error_reason = "user.age_not_qualified"
    variables = {
        verification_result = {
            status = "failed",
            message = "用户年龄未满18岁",
            age_info = {
                current_age = userAge,
                required_age = 18
            }
        }
    }
    return
end

-- 5. 身份证号验证（必须提供有效的身份证号）
local personID = user.personId
if not personID or personID == "" then
    valid = false
    action = "require_id_card"
    error = "用户未提供身份证号"
    error_reason = "user.id_card_required"
    variables = {
        verification_result = {
            status = "failed",
            message = "用户未提供身份证号",
            required_action = "provide_id_card"
        }
    }
    return
end

-- 6. 身份证号格式验证（简单验证18位）
local idCardLength = fast_contains(personID, "") and #personID or 0
if idCardLength ~= 18 then
    valid = false
    action = "reject"
    error = "身份证号格式不正确"
    error_reason = "user.id_card_format_invalid"
    variables = {
        verification_result = {
            status = "failed",
            message = "身份证号格式不正确",
            id_card_info = {
                person_id = personID,
                length = idCardLength,
                required_length = 18
            }
        }
    }
    return
end

-- 7. 手机号验证（必须提供有效的手机号）
local phone = user.phone
if not phone or phone == "" then
    valid = false
    action = "require_phone"
    error = "用户未提供手机号"
    error_reason = "user.phone_required"
    variables = {
        verification_result = {
            status = "failed",
            message = "用户未提供手机号",
            required_action = "provide_phone"
        }
    }
    return
end

-- 8. 手机号格式验证（简单验证11位数字）
local phoneLength = fast_contains(phone, "") and #phone or 0
if phoneLength ~= 11 then
    valid = false
    action = "reject"
    error = "手机号格式不正确"
    error_reason = "user.phone_format_invalid"
    variables = {
        verification_result = {
            status = "failed",
            message = "手机号格式不正确",
            phone_info = {
                phone = phone,
                length = phoneLength,
                required_length = 11
            }
        }
    }
    return
end

-- 9. 用户等级验证（VIP用户优先处理）
local userLevel = user.level
local isVIP = userLevel >= 1 and userLevel <= 4

-- 10. 测试账号验证（测试账号不能进行实名操作）
local isTestAccount = user.testAccount == 1
if isTestAccount then
    valid = false
    action = "reject"
    error = "测试账号不能进行实名操作"
    error_reason = "user.test_account_not_allowed"
    variables = {
        verification_result = {
            status = "failed",
            message = "测试账号不能进行实名操作",
            account_info = {
                test_account = user.testAccount
            }
        }
    }
    return
end

-- 所有验证通过
valid = true
action = "approve"

variables = {
    verification_result = {
        status = "success",
        message = "用户实名验证通过",
        user_info = {
            id = user.id,
            nickname = user.nickname,
            phone = user.phone,
            account = user.account,
            email = user.email,
            age = user.age,
            gender = user.gender,
            person_id = user.personId,
            is_real = user.isReal,
            real_at = user.realAt,
            level = user.level,
            test_account = user.testAccount,
            status = user.status
        },
        verification_details = {
            real_name_verified = true,
            age_qualified = true,
            id_card_provided = true,
            phone_provided = true,
            user_status_valid = true,
            test_account_check = false,
            is_vip = isVIP,
            user_level = userLevel
        },
        timestamp = now()
    }
} 