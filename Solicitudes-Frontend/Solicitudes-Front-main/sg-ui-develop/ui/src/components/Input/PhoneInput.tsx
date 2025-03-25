import { useState } from "react";

type CountryInfo = {
  name: string;
  code: string;
  flag: JSX.Element;
};

const countries: CountryInfo[] = [
  {
    name: "Argentina",
    code: "+54",
    flag: (
      <svg className="h-4 w-4 me-2" fill="none" viewBox="0 0 20 15">
        <rect width="19.6" height="14" y=".5" fill="#fff" rx="2" />
        <mask
          id="a"
          style={{ maskType: "luminance" }}
          width="20"
          height="15"
          x="0"
          y="0"
          maskUnits="userSpaceOnUse"
        >
          <rect width="19.6" height="14" y=".5" fill="#fff" rx="2" />
        </mask>
        <g mask="url(#a)">
          <path fill="#fff" d="M0 .5h19.6v14H0z" />
          <path fill="#88BBE8" d="M0 .5h19.6v4.667H0zm0 9.333h19.6V14.5H0z" />
          <path
            fill="#F4B32E"
            d="M9.8 6.1l-.766.533.228-.888-.624-.655.836-.063.326-.851.326.851.836.063-.624.655.228.888z"
          />
        </g>
      </svg>
    ),
  },
  {
    name: "Brasil",
    code: "+55",
    flag: (
      <svg className="h-4 w-4 me-2" fill="none" viewBox="0 0 20 15">
        <rect width="19.6" height="14" y=".5" fill="#fff" rx="2" />
        <mask
          id="a"
          style={{ maskType: "luminance" }}
          width="20"
          height="15"
          x="0"
          y="0"
          maskUnits="userSpaceOnUse"
        >
          <rect width="19.6" height="14" y=".5" fill="#fff" rx="2" />
        </mask>
        <g mask="url(#a)">
          <path fill="#229E45" d="M0 .5h19.6v14H0z" />
          <path
            fill="#F8E509"
            d="M2.333 7.5l7.467-4.667 7.467 4.667-7.467 4.667z"
          />
          <circle cx="9.8" cy="7.5" r="2.8" fill="#2B49A3" />
        </g>
      </svg>
    ),
  },
  {
    name: "Chile",
    code: "+56",
    flag: (
      <svg className="h-4 w-4 me-2" fill="none" viewBox="0 0 20 15">
        <rect width="19.6" height="14" y=".5" fill="#fff" rx="2" />
        <mask
          id="a"
          style={{ maskType: "luminance" }}
          width="20"
          height="15"
          x="0"
          y="0"
          maskUnits="userSpaceOnUse"
        >
          <rect width="19.6" height="14" y=".5" fill="#fff" rx="2" />
        </mask>
        <g mask="url(#a)">
          <path fill="#fff" d="M0 .5h19.6v7H0z" />
          <path fill="#EA3B2E" d="M0 7.5h19.6v7H0z" />
          <path fill="#0B48C2" d="M0 .5h6.533v7H0z" />
          <path
            fill="#fff"
            d="M3.267 2.367l-.766.533.228-.888-.624-.655.836-.063.326-.851.326.851.836.063-.624.655.228.888z"
          />
        </g>
      </svg>
    ),
  },
  {
    name: "Colombia",
    code: "+57",
    flag: (
      <svg className="h-4 w-4 me-2" fill="none" viewBox="0 0 20 15">
        <rect width="19.6" height="14" y=".5" fill="#fff" rx="2" />
        <mask
          id="a"
          style={{ maskType: "luminance" }}
          width="20"
          height="15"
          x="0"
          y="0"
          maskUnits="userSpaceOnUse"
        >
          <rect width="19.6" height="14" y=".5" fill="#fff" rx="2" />
        </mask>
        <g mask="url(#a)">
          <path fill="#FFD935" d="M0 .5h19.6v7H0z" />
          <path fill="#0748AE" d="M0 7.5h19.6v3.5H0z" />
          <path fill="#DE2035" d="M0 11h19.6v3.5H0z" />
        </g>
      </svg>
    ),
  },
  {
    name: "México",
    code: "+52",
    flag: (
      <svg className="h-4 w-4 me-2" fill="none" viewBox="0 0 20 15">
        <rect width="19.6" height="14" y=".5" fill="#fff" rx="2" />
        <mask
          id="a"
          style={{ maskType: "luminance" }}
          width="20"
          height="15"
          x="0"
          y="0"
          maskUnits="userSpaceOnUse"
        >
          <rect width="19.6" height="14" y=".5" fill="#fff" rx="2" />
        </mask>
        <g mask="url(#a)">
          <path fill="#fff" d="M6.533.5h6.533v14H6.533z" />
          <path fill="#E3283E" d="M13.067.5H19.6v14h-6.533z" />
          <path fill="#128A60" d="M0 .5h6.533v14H0z" />
          <circle cx="9.8" cy="7.5" r="1.867" fill="#8C9157" />
        </g>
      </svg>
    ),
  },
  {
    name: "España",
    code: "+34",
    flag: (
      <svg className="h-4 w-4 me-2" fill="none" viewBox="0 0 20 15">
        <rect width="19.6" height="14" y=".5" fill="#fff" rx="2" />
        <mask
          id="a"
          style={{ maskType: "luminance" }}
          width="20"
          height="15"
          x="0"
          y="0"
          maskUnits="userSpaceOnUse"
        >
          <rect width="19.6" height="14" y=".5" fill="#fff" rx="2" />
        </mask>
        <g mask="url(#a)">
          <path fill="#DD172C" d="M0 .5h19.6v14H0z" />
          <path fill="#FFD133" d="M0 3.833h19.6v7.333H0z" />
        </g>
      </svg>
    ),
  },
  {
    name: "United States",
    code: "+1",
    flag: (
      <svg className="h-4 w-4 me-2" fill="none" viewBox="0 0 20 15">
        <rect width="19.6" height="14" y=".5" fill="#fff" rx="2" />
        <mask
          id="a"
          style={{ maskType: "luminance" }}
          width="20"
          height="15"
          x="0"
          y="0"
          maskUnits="userSpaceOnUse"
        >
          <rect width="19.6" height="14" y=".5" fill="#fff" rx="2" />
        </mask>
        <g mask="url(#a)">
          <path
            fill="#D02F44"
            fillRule="evenodd"
            d="M19.6.5H0v.933h19.6V.5zm0 1.867H0V3.3h19.6v-.933zM0 4.233h19.6v.934H0v-.934zM19.6 6.1H0v.933h19.6V6.1zM0 7.967h19.6V8.9H0v-.933zm19.6 1.866H0v.934h19.6v-.934zM0 11.7h19.6v.933H0V11.7zm19.6 1.867H0v.933h19.6v-.933z"
            clipRule="evenodd"
          />
          <path fill="#46467F" d="M0 .5h8.4v6.533H0z" />
        </g>
      </svg>
    ),
  },
];

interface PhoneInputProps {
  phone: string;
  setPhone: (value: string) => void;
}

const PhoneInput: React.FC<PhoneInputProps> = ({ phone, setPhone }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [selectedCountry, setSelectedCountry] = useState(countries[0]);

  const handleCountrySelect = (country: CountryInfo) => {
    setSelectedCountry(country);
    setIsOpen(false);
  };

  const formatPhoneNumber = (value: string) => {
    const phone = value.replace(/\D/g, "");
    if (phone.length < 4) return phone;
    if (phone.length < 7) return `${phone.slice(0, 3)}-${phone.slice(3)}`;
    return `${phone.slice(0, 3)}-${phone.slice(3, 6)}-${phone.slice(6, 10)}`;
  };

  const handlePhoneChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const formatted = formatPhoneNumber(e.target.value);
    setPhone(formatted);
  };

  return (
    <div className="flex items-center">
      <div className="relative">
        <button
          type="button"
          onClick={() => setIsOpen(!isOpen)}
          className="flex-shrink-0 z-10 inline-flex items-center py-2.5 px-4 text-sm font-medium text-gray-900 bg-gray-100 border border-gray-300 rounded-s-lg hover:bg-gray-200 focus:ring-4 focus:outline-none focus:ring-gray-100"
        >
          {selectedCountry.flag}
          {selectedCountry.code}
        </button>

        {isOpen && (
          <div className="absolute mt-2 z-10 w-52 bg-white divide-y divide-gray-100 rounded-lg shadow dark:bg-gray-700">
            <ul className="py-2 text-sm text-gray-700 dark:text-gray-200">
              {countries.map((country) => (
                <li key={country.code}>
                  <button
                    type="button"
                    className="inline-flex w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-600 dark:hover:text-white"
                    onClick={() => handleCountrySelect(country)}
                  >
                    <span className="inline-flex items-center">
                      {country.flag}
                      {country.name} ({country.code})
                    </span>
                  </button>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>

      <div className="relative w-full">
        <input
          type="tel"
          value={phone}
          onChange={handlePhoneChange}
          className="block p-2.5 w-full z-20 text-sm text-gray-900 bg-gray-50 rounded-e-lg border-s-0 border border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-s-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:border-blue-500"
          placeholder="123-456-7890"
          pattern="[0-9]{3}-[0-9]{3}-[0-9]{4}"
          required
        />
      </div>
    </div>
  );
};

export default PhoneInput;
