import { Button } from "@/components/ui/button";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { cn } from "@/lib/utils";
import BasicTableFilter from "@/components/Table/BasicTableFilter.tsx";
import { getStatusOptions } from "@/components/Table/statusOptions.ts";

interface BasicTableProps {
    headers: Array<{
        id: string,
        title: string,
        filter?: {
            typeOptions: string[]
            selectedType: string | null
            setSelectedType: (type: string | null) => void
        }
    }>;
    data: any[];
    actions: Array<{
        label: string,
        displayStatus: string[],
        callback: (data: any) => void
    }>;
}

export default function BasicTable({ headers, data, actions }: BasicTableProps) {
    return (
        <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
            <Table>
                <TableHeader>
                    <TableRow className="hover:bg-transparent">
                        {
                            headers.map((header) => (
                                <TableHead key={ header.title } className="bg-gray-100 text-black font-medium">
                                    <div className="flex items-center gap-1">
                                        { header.title }
                                        {
                                            header.filter &&
                                            <BasicTableFilter
                                                typeOptions={ header.filter.typeOptions }
                                                selectedType={ header.filter.selectedType }
                                                setSelectedType={ header.filter.setSelectedType }
                                            ></BasicTableFilter>
                                        }
                                    </div>
                                </TableHead>
                            ))
                        }
                    </TableRow>
                </TableHeader>
                <TableBody>
                    { data.length === 0 ? (
                        <TableRow>
                            <TableCell
                                colSpan={ 6 }
                                className="h-24 text-center text-gray-500"
                            >
                                Por el momento, no tenés acciones para realizar.
                            </TableCell>
                        </TableRow>
                    ) : (
                        data.map((element, index: number) => (
                            <TableRow key={ index } className="hover:bg-gray-50 h-16">
                                {
                                    headers.map((header) => {
                                        switch (header.id) {
                                            case "status":
                                                return <TableCell className="py-4">
                                                    <span
                                                        className={ cn(
                                                            "px-3 py-1 rounded-full text-sm font-normal text-white",
                                                            getStatusOptions(element[header.id]).styles,
                                                        ) }
                                                    >
                                                      { getStatusOptions(element[header.id]).label }
                                                    </span>
                                                </TableCell>;
                                            case "actions":
                                                return <TableCell className="py-4 text-right">
                                                    {
                                                        actions.reduce((acc, action) => {
                                                            if (action.displayStatus.includes(element.status)) {
                                                                acc = (
                                                                    <Button
                                                                        variant="default"
                                                                        className="px-6 bg-[#526D4E] hover:bg-[#526D4E]/90"
                                                                        onClick={ () => action.callback(element) }
                                                                    >
                                                                        { action.label }
                                                                    </Button>
                                                                )
                                                            }
                                                            return acc
                                                        }, <></>)
                                                    }
                                                </TableCell>;
                                            case "requesterAddress":
                                                return <TableCell className="py-4">{ element[header.id] }</TableCell>;
                                            default:
                                                return <TableCell className="py-4 text-nowrap">{ element[header.id] }</TableCell>;
                                        }
                                    })
                                }
                            </TableRow>
                        ))
                    ) }
                </TableBody>
            </Table>
        </div>
    );
}
