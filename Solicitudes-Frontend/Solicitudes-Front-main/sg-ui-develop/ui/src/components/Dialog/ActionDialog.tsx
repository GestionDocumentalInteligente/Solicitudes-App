import { AlertTriangle, type LucideIcon } from 'lucide-react';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog.tsx";
import { Button } from "@/components/ui/button.tsx";
import { ReactNode } from "react";
import { cn } from "@/lib/utils.ts";

type ActionDialogProps = {
    title: string,
    description: string,
    icon?: LucideIcon;
    showDialog: boolean
    resetDialog: () => void
    children?: ReactNode;
    primaryButtonLabel: string;
    primaryButtonAction: () => void;
    secondaryButtonLabel: string;
}

export default function ActionDialog({
                                        title,
                                        description,
                                        icon,
                                        showDialog,
                                        resetDialog,
                                        children,
                                        primaryButtonLabel,
                                        primaryButtonAction,
                                        secondaryButtonLabel,
                                     }: ActionDialogProps) {
    const Icon = icon || AlertTriangle;
    return (
        <Dialog open={ showDialog }>
            <DialogContent className={ cn(children ? "max-w-3xl" : "max-w-xl" )}>
                <DialogHeader>
                    {
                        icon && <Icon className="w-16 h-16 stroke-1 m-auto"></Icon>
                    }
                    <DialogTitle>{ title }</DialogTitle>
                    <DialogDescription>{ description }</DialogDescription>
                </DialogHeader>

                { children }

                <div className="flex justify-end gap-4">
                    <Button
                        variant="outline"
                        className="w-1/2"
                        onClick={ () => resetDialog() }
                    >
                        { secondaryButtonLabel }
                    </Button>
                    <Button
                        className="bg-[#3b5c3f] hover:bg-[#2e4831] text-white w-1/2"
                        onClick={ () => {
                            primaryButtonAction();
                            resetDialog();
                        } }
                    >
                        { primaryButtonLabel }
                    </Button>
                </div>
            </DialogContent>
        </Dialog>
    );
}