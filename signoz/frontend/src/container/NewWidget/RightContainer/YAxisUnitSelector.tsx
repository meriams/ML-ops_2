import { AutoComplete, Col, Input, Typography } from 'antd';
import { find } from 'lodash-es';
import { Dispatch, SetStateAction } from 'react';

import { flattenedCategories } from './dataFormatCategories';

const findCategoryById = (
	searchValue: string,
): Record<string, string> | undefined =>
	find(flattenedCategories, (option) => option.id === searchValue);
const findCategoryByName = (
	searchValue: string,
): Record<string, string> | undefined =>
	find(flattenedCategories, (option) => option.name === searchValue);

function YAxisUnitSelector({
	defaultValue,
	onSelect,
	fieldLabel,
}: {
	defaultValue: string;
	onSelect: Dispatch<SetStateAction<string>>;
	fieldLabel: string;
}): JSX.Element {
	const onSelectHandler = (selectedValue: string): void => {
		onSelect(findCategoryByName(selectedValue)?.id || '');
	};
	const options = flattenedCategories.map((options) => ({
		value: options.name,
	}));
	return (
		<Col style={{ marginBottom: 12, marginTop: 12 }}>
			<Typography.Text>{fieldLabel}</Typography.Text>
			<AutoComplete
				style={{ width: '100%' }}
				options={options}
				defaultValue={findCategoryById(defaultValue)?.name}
				onSelect={onSelectHandler}
				filterOption={(inputValue, option): boolean => {
					if (option) {
						return (
							option.value.toUpperCase().indexOf(inputValue.toUpperCase()) !== -1
						);
					}
					return false;
				}}
			>
				<Input size="large" placeholder="Unit" allowClear />
			</AutoComplete>
		</Col>
	);
}

export default YAxisUnitSelector;
