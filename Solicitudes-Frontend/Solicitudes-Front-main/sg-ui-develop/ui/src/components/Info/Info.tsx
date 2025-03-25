interface InfoProps {
  text: string | React.ReactNode;
}

const Info: React.FC<InfoProps> = ({ text }) => {
  return <p className="mt-2 text-sm text-gray-500">{text}</p>;
};

export default Info;
