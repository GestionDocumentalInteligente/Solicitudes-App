import { create } from 'zustand';
import { type LucideIcon } from 'lucide-react';
import { ReactNode } from "react";

interface ActionDialogState {
    isVisible: boolean;
    title: string;
    description: string;
    icon?: LucideIcon;
    children?: ReactNode;
    primaryButtonLabel: string;
    primaryButtonAction: () => void;
    secondaryButtonLabel: string;
    showDialog: (params: Omit<ActionDialogState, 'isVisible' | 'showDialog' | 'hideDialog'>) => void;
    hideDialog: () => void;
}

export const useActionDialogStore = create<ActionDialogState>((set) => ({
    isVisible: false,
    title: '',
    description: '',
    primaryButtonLabel: '',
    primaryButtonAction: () => {
    },
    secondaryButtonLabel: '',
    showDialog: (params) => set({ ...params, isVisible: true }),
    hideDialog: () => set({
        isVisible: false, title: '', description: '', primaryButtonLabel: '', primaryButtonAction: () => {
        }, secondaryButtonLabel: '',
    }),
}));