interface LabelProps {
  text: string;
  mandatory: boolean;
}

export const Label: React.FC<LabelProps> = ({ text, mandatory }) => {
  return (
    <label htmlFor="address" className="block text-sm font-bold text-gray-700">
      {text} {mandatory && <span className="text-red-500">*</span>}
    </label>
  );
};
