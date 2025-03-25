import { create } from 'zustand'
import { type LucideIcon } from 'lucide-react'

interface StatusMessageState {
    isVisible: boolean
    type: 'success' | 'error'
    icon?: LucideIcon
    title: string[]
    description: string[]
    buttonText: string
    onAction?: () => void
    showMessage: (params: Omit<StatusMessageState, 'isVisible' | 'showMessage' | 'hideMessage'>) => void
    hideMessage: () => void
}

export const useStatusMessageStore = create<StatusMessageState>((set) => ({
    isVisible: false,
    type: 'success',
    title: [''],
    description: [''],
    buttonText: '',
    showMessage: (params) => set({ ...params, isVisible: true }),
    hideMessage: () => set({ isVisible: false, type: 'success', title: [''], description: [''], buttonText: '' })
}))