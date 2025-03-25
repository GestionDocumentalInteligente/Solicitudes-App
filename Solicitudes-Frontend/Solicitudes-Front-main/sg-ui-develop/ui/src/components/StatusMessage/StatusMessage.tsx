import * as React from "react";
import { Button } from "@/components/ui/button.tsx";
import { type LucideIcon, FileX, CircleCheck } from 'lucide-react';

type StatusProps = {
    type: 'success' | 'error'
    icon?: LucideIcon
    title: string[]
    description: string[]
    buttonText: string
    onAction: () => void
}


export default function StatusMessage({
                                          type = 'success',
                                          icon,
                                          title,
                                          description,
                                          buttonText,
                                          onAction,
                                      }: StatusProps) {
    const Icon = icon || (type === 'success' ? CircleCheck : FileX);

    const renderText = (text: string[]) => {
        return text.map((line, index) => (
            <React.Fragment key={ index }>
                { line }
                { index < title.length - 1 && <br/> }
            </React.Fragment>
        ));
    };

    return (
        <div className="min-h-screen bg-white flex items-center justify-center">
            <div className="w-full mx-auto px-4">
                <div className="flex flex-col items-center text-center space-y-6">
                    <Icon className="w-16 h-16 stroke-1"/>
                    <div className="space-y-4">
                        <h1 className="text-3xl font-bold">{ renderText(title) }</h1>
                        <div className="space-y-2">
                            <p className="text-xl">{ renderText(description) }</p>
                        </div>
                    </div>
                    <Button
                        className="w-full max-w-md bg-[#3b5c3f] hover:bg-[#2e4831] text-white"
                        onClick={ onAction }
                    >
                        { buttonText }
                    </Button>
                </div>
            </div>
        </div>
    );
}