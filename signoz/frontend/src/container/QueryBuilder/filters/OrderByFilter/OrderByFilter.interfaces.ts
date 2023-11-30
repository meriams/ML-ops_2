import {
	IBuilderQuery,
	OrderByPayload,
} from 'types/api/queryBuilder/queryBuilderData';

export type OrderByFilterProps = {
	query: IBuilderQuery;
	onChange: (values: OrderByPayload[]) => void;
};

export type OrderByFilterValue = {
	disabled?: boolean;
	key: string;
	label: string;
	title?: string;
	value: string;
};
