export interface Validation {
    recordNumber: string;
    requestType: string;
    documentType: string;
    deliveryDate: string;
    status: string;
    requesterFullName: string;
    requesterCuil: string;
    requesterAddress: string;
    grantedToAuth?: string;
    requiresInsurance: boolean;
    estimatedTime: string;
    assignedTask: string[];
    replacementIFs: Array<{
        id: number;
        title: string;
        gedoCode: string;
    }>
    documents: Document[];
}

export interface Document {
    id: string
    title: string,
    base64?: string,
    gedoCode?: string,
    verifiedBy?: string,
    verifiedDate?: string,
    sections: Section[]
}

export interface Section {
    id: string;
    title: string;
    status: ValidationStatus;
    question?: string;
    questionClarification?: string;
    observation?: string;
    observationPlaceholder?: string;
    observationOptions?: Array<{
        id: string;
        description: string;
    }>;
    selectedObservationOptions?: string[];
    isConfirmed?: boolean;
    inputType?: InputType;
    inputPlaceholder?: string;
    inputOptions?: string[];
    value?: string;
    isExpanded: boolean;
    show: boolean;
}

export type ValidationStatus = 'valid' | 'invalid' | null

type InputType = 'text' | 'radio' | null


// TODO: GENERALIZE
export interface RequestResponse {
    data: Validation[];
    error: string;
}

class RequestError extends Error {
}

class UnhandledError extends Error {
}

export { RequestError, UnhandledError };