import { CardInfo } from "./types";
import { Link } from "react-router-dom";

interface CardProps {
  info: CardInfo;
  icon: JSX.Element | null;
  to: string;
}

export function Card({ info, icon, to }: CardProps) {
  return (
    <Link
      to={to}
      key={info.id}
      className={`card p-6 rounded-lg shadow-lg transition flex items-center max-h-36 overflow-hidden ${
        info.is_active
          ? "bg-white border-2 border-gray-400 cursor-pointer hover:shadow-xl"
          : "bg-gray-200 border-2 border-gray-400 cursor-not-allowed"
      }`}
    >
      {icon && <div className="flex-shrink-0 text-center mr-4">{icon}</div>}
      <div className="flex flex-col">
        <p
          className={`text-base font-base ${
            info.is_active ? "text-black" : "text-gray-500"
          }`}
        >
          {info.name}
        </p>

        {info.description && (
          <p
            className={`text-sm ${
              info.is_active ? "text-black" : "text-gray-500"
            } line-clamp-3 transition-all duration-300`}
          >
            {info.description}
          </p>
        )}
      </div>
    </Link>
  );
}
