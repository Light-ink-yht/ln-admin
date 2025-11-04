import { request } from '@/utils/request'
import type { ResponseData } from '@/utils/request'

/**
 * 短信模板信息
 */
export interface SMSTemplate {
    templateId: string
    type: string
    templateName: string
    tencentTemplateId: string
    title: string
    content: string
    description?: string
    status?: string
    creatorId?: string
    modifierId?: string
    createdAt?: string
    updatedAt?: string
}

/**
 * 短信模板列表请求参数
 */
export interface SMSTemplateListParams {
    page?: number
    pageSize?: number
    type?: string
    status?: string
}

/**
 * 短信验证码信息
 */
export interface SMSCode {
    codeId: string
    phone: string
    code: string
    type: string
    status: string
    expireAt?: string
    usedAt?: string
    ip: string
    sendCount: number
    createdAt: string
}

/**
 * 短信验证码列表请求参数
 */
export interface SMSCodeListParams {
    page?: number
    pageSize?: number
    phone?: string
    type?: string
    status?: string
}

/**
 * 创建短信模板请求
 */
export interface CreateSMSTemplateRequest {
    type: string
    templateName: string
    tencentTemplateId: string
    title: string
    content: string
    description?: string
    status?: string
}

/**
 * 更新短信模板请求
 */
export interface UpdateSMSTemplateRequest {
    templateName?: string
    tencentTemplateId?: string
    title?: string
    content?: string
    description?: string
    status?: string
}

/**
 * 短信 API
 */
export const smsApi = {
    /**
     * 获取短信模板列表
     */
    getTemplateList(params: SMSTemplateListParams): Promise<ResponseData<{ list: SMSTemplate[]; total: number }>> {
        return request.get<{ list: SMSTemplate[]; total: number }>('/sms/template/list', params)
    },

    /**
     * 获取短信模板详情
     */
    getTemplateDetail(templateId: string): Promise<ResponseData<SMSTemplate>> {
        return request.get<SMSTemplate>(`/sms/template/${templateId}`)
    },

    /**
     * 创建短信模板
     */
    createTemplate(data: CreateSMSTemplateRequest): Promise<ResponseData<void>> {
        return request.post<void>('/sms/template', data)
    },

    /**
     * 更新短信模板
     */
    updateTemplate(templateId: string, data: UpdateSMSTemplateRequest): Promise<ResponseData<void>> {
        return request.put<void>(`/sms/template/${templateId}`, data)
    },

    /**
     * 删除短信模板
     */
    deleteTemplate(templateId: string): Promise<ResponseData<void>> {
        return request.delete<void>(`/sms/template/${templateId}`)
    },

    /**
     * 获取短信验证码列表
     */
    getCodeList(params: SMSCodeListParams): Promise<ResponseData<{ list: SMSCode[]; total: number }>> {
        return request.get<{ list: SMSCode[]; total: number }>('/sms/code/list', params)
    },
}

