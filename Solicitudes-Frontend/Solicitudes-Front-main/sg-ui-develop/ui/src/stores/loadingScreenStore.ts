import { create } from 'zustand';

interface LoadingScreenState {
    isVisible: boolean;
    title?: string[];
    description?: string[];
    showLoading: (params?: Omit<LoadingScreenState, 'isVisible' | 'showLoading' | 'hideLoading'>) => void;
    hideLoading: () => void;
}

export const useLoadingScreenStore = create<LoadingScreenState>((set) => ({
    isVisible: false,
    title: [''],
    description: [''],
    showLoading: (params) => set({ ...params, isVisible: true }),
    hideLoading: () => set({ isVisible: false, title: [''], description: [''] }),
}));