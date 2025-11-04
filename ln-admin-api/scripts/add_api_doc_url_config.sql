-- 添加接口文档地址配置项
-- 如果配置已存在则更新，否则插入新配置

INSERT INTO AA02 (
    AAA001,  -- ConfigKey
    AAA002,  -- ConfigValue
    AAA003,  -- ConfigName
    AAA004,  -- ConfigGroup
    AAA005,  -- Description
    AAA006,  -- Status
    AAA007,  -- CreatorID
    AAA008,  -- ModifierID
    AAA009,  -- CreatedAt
    AAA010   -- UpdatedAt
) VALUES (
    'api_doc_url',
    '',
    '接口文档地址',
    'system',
    'API接口文档地址，用于在接口文档页面进行跳转',
    '1',
    'system',
    'system',
    NOW(),
    NOW()
)
ON DUPLICATE KEY UPDATE
    AAA003 = '接口文档地址',
    AAA004 = 'system',
    AAA005 = 'API接口文档地址，用于在接口文档页面进行跳转',
    AAA008 = 'system',
    AAA010 = NOW();

