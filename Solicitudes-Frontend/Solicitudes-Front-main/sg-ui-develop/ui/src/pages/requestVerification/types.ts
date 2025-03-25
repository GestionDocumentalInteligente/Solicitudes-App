export interface Verification {
    verificationCase: string;
    invokedCharacter: string;
    recordNumber: string;
    requestType: string;
    documentType: string;
    deliveryDate: string;
    status: string;
    status_task: string;
    status_property: string;
    requesterFullName: string;
    requesterCuil: string;
    requesterAddress: string;
    documents: Document[];
}

export interface Document {
    id: string
    title: string,
    base64?: string,
    gedoCode?: string,
    request: DocumentContent
    sections: Section[]
}

export interface DocumentContent {
    descriptiveMemory: string
    assignedTask: string[]
    estimatedTime: string
    requiresInsurance: string
}

export interface Section {
    id: string;
    title: string;
    status: VerificationStatus;
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

export type VerificationStatus = 'valid' | 'invalid' | null

type InputType = 'text' | 'radio' | null


// TODO: GENERALIZE
export interface RequestResponse {
    data: Verification[];
    error: string;
}

class RequestError extends Error {
}

class UnhandledError extends Error {
}

export { RequestError, UnhandledError };