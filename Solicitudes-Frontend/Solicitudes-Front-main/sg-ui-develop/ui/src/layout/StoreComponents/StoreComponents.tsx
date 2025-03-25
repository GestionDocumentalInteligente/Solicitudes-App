import StatusMessage from "@/components/StatusMessage/StatusMessage.tsx";
import LoadingScreen from "@/components/LoadingScreen/LoadingScreen.tsx";
import ActionDialog from "@/components/Dialog/ActionDialog.tsx";
import { useStatusMessageStore } from "@/stores/statusMessageStore.ts";
import { useLoadingScreenStore } from "@/stores/loadingScreenStore.ts";
import { useActionDialogStore } from "@/stores/actionDialogStore.ts";

export default function StoreComponents() {
    const statusMessage = useStatusMessageStore();
    const loadingScreen = useLoadingScreenStore();
    const actionDialog = useActionDialogStore();

    const renderStatusMessage = () => {
        if (!statusMessage.isVisible) return null;

        return (
            <div className="fixed inset-0 z-50">
                <StatusMessage
                    type={statusMessage.type}
                    icon={statusMessage.icon}
                    title={statusMessage.title}
                    description={statusMessage.description}
                    buttonText={statusMessage.buttonText}
                    onAction={() => {
                        statusMessage.onAction?.();
                        statusMessage.hideMessage();
                    }}
                />
            </div>
        );
    };

    const renderLoadingScreen = () => {
        if (!loadingScreen.isVisible) return null;

        return (
            <div className="fixed inset-0 z-20">
                <LoadingScreen
                    title={loadingScreen.title}
                    description={loadingScreen.description}
                ></LoadingScreen>
            </div>
        );
    };

    const renderActionDialog = () => {
        // TODO: CHECK IF CONDITION IS NEEDED
        if (!actionDialog.isVisible) return null;

        return (
            <div className="fixed inset-0 z-50">
                <ActionDialog
                    title={actionDialog.title}
                    description={actionDialog.description}
                    showDialog={actionDialog.isVisible}
                    resetDialog={actionDialog.hideDialog}
                    icon={actionDialog.icon}
                    primaryButtonLabel={actionDialog.primaryButtonLabel}
                    primaryButtonAction={actionDialog.primaryButtonAction}
                    secondaryButtonLabel={actionDialog.secondaryButtonLabel}
                >
                    {actionDialog.children}
                </ActionDialog>
            </div>
        );
    }

    return (
        <>
            {renderLoadingScreen()}
            {renderStatusMessage()}
            {renderActionDialog()}
        </>
    );
}