import { Banner, Tooltip } from "flowbite-react";
import { FaRegCircleQuestion } from "react-icons/fa6";

type Placement =
  | "top"
  | "right"
  | "bottom"
  | "left"
  | "top-start"
  | "top-end"
  | "right-start"
  | "right-end"
  | "bottom-start"
  | "bottom-end"
  | "left-start"
  | "left-end";

interface TooltipContentProps {
  title: string;
  content: React.ReactNode;
}

const TooltipContent = ({ title, content }: TooltipContentProps) => {
  return (
    <Banner className="max-w-xs break-words z-70">
      <h2 className="mb-1 text-base font-bold text-gray-900">{title}</h2>
      {content}
    </Banner>
  );
};

interface FormInfoProps {
  info: string;
  placement: Placement;
  tooltipContent: TooltipContentProps;
}

export const FormInfo = ({
  info,
  placement,
  tooltipContent,
}: FormInfoProps) => {
  return (
    <div className="flex items-center text-sm font-normal text-gray-500">
      <span className="[&_p]:inline m-2 font-semibold">{info}</span>
      <div className="flex gap-2">
        <Tooltip
          style="light"
          placement={placement}
          content={
            <TooltipContent
              title={tooltipContent.title}
              content={tooltipContent.content}
            />
          }
          arrow={false}
        >
          <FaRegCircleQuestion className="mr-4 h-4 w-4 z-70" />
        </Tooltip>
      </div>
    </div>
  );
};
