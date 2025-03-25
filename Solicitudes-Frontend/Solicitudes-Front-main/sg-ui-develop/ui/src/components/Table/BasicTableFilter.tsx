import { Button } from "@/components/ui/button";
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Label } from "@/components/ui/label";
import { ChevronDown } from 'lucide-react';

interface FilterProps {
    typeOptions: string[];
    selectedType: string | null;
    setSelectedType: (type: string | null) => void;
}

export default function BasicTableFilter({ typeOptions, selectedType, setSelectedType }: FilterProps) {
    return (
        <Popover>
            <PopoverTrigger asChild>
                <Button variant="ghost" size="sm"
                        className="h-4 w-4 p-0 hover:bg-transparent">
                    <ChevronDown className="h-4 w-4 text-gray-900"/>
                </Button>
            </PopoverTrigger>
            <PopoverContent className="w-48">
                <RadioGroup
                    value={ selectedType || "" }
                    onValueChange={ setSelectedType }
                >
                    { typeOptions.map((option) => (
                        <div key={ option } className="flex items-center space-x-2">
                            <RadioGroupItem value={ option }
                                            id={ `type-${ option }` }/>
                            <Label htmlFor={ `type-${ option }` }>{ option }</Label>
                        </div>
                    )) }
                </RadioGroup>
                <Button
                    variant="ghost"
                    size="sm"
                    onClick={ () => setSelectedType(null) }
                    className="w-full mt-2"
                >
                    Limpiar filtro
                </Button>
            </PopoverContent>
        </Popover>
    );
}
