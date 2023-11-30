/* eslint-disable sonarjs/cognitive-complexity */
import { orange } from '@ant-design/colors';
import {
	Button,
	Col,
	Divider,
	Input,
	Select,
	Switch,
	Tag,
	Typography,
} from 'antd';
import query from 'api/dashboard/variables/query';
import Editor from 'components/Editor';
import { commaValuesParser } from 'lib/dashbaordVariables/customCommaValuesParser';
import sortValues from 'lib/dashbaordVariables/sortVariableValues';
import { map } from 'lodash-es';
import { useEffect, useState } from 'react';
import {
	IDashboardVariable,
	TSortVariableValuesType,
	TVariableQueryType,
	VariableQueryTypeArr,
	VariableSortTypeArr,
} from 'types/api/dashboard/getAll';
import { v4 } from 'uuid';

import { variablePropsToPayloadVariables } from '../../../utils';
import { TVariableViewMode } from '../types';
import { LabelContainer, VariableItemRow } from './styles';

const { Option } = Select;

interface VariableItemProps {
	variableData: IDashboardVariable;
	existingVariables: Record<string, IDashboardVariable>;
	onCancel: () => void;
	onSave: (name: string, arg0: IDashboardVariable, arg1: string) => void;
	validateName: (arg0: string) => boolean;
	variableViewMode: TVariableViewMode;
}
function VariableItem({
	variableData,
	existingVariables,
	onCancel,
	onSave,
	validateName,
	variableViewMode,
}: VariableItemProps): JSX.Element {
	const [variableName, setVariableName] = useState<string>(
		variableData.name || '',
	);
	const [variableDescription, setVariableDescription] = useState<string>(
		variableData.description || '',
	);
	const [queryType, setQueryType] = useState<TVariableQueryType>(
		variableData.type || 'QUERY',
	);
	const [variableQueryValue, setVariableQueryValue] = useState<string>(
		variableData.queryValue || '',
	);
	const [variableCustomValue, setVariableCustomValue] = useState<string>(
		variableData.customValue || '',
	);
	const [variableTextboxValue, setVariableTextboxValue] = useState<string>(
		variableData.textboxValue || '',
	);
	const [
		variableSortType,
		setVariableSortType,
	] = useState<TSortVariableValuesType>(
		variableData.sort || VariableSortTypeArr[0],
	);
	const [variableMultiSelect, setVariableMultiSelect] = useState<boolean>(
		variableData.multiSelect || false,
	);
	const [variableShowALLOption, setVariableShowALLOption] = useState<boolean>(
		variableData.showALLOption || false,
	);
	const [previewValues, setPreviewValues] = useState<string[]>([]);

	// Internal states
	const [previewLoading, setPreviewLoading] = useState<boolean>(false);
	// Error messages
	const [errorName, setErrorName] = useState<boolean>(false);
	const [errorPreview, setErrorPreview] = useState<string | null>(null);

	useEffect(() => {
		setPreviewValues([]);
		if (queryType === 'CUSTOM') {
			setPreviewValues(
				sortValues(
					commaValuesParser(variableCustomValue),
					variableSortType,
				) as never,
			);
		}
	}, [
		queryType,
		variableCustomValue,
		variableData.customValue,
		variableData.type,
		variableSortType,
	]);

	const handleSave = (): void => {
		const newVariableData: IDashboardVariable = {
			name: variableName,
			description: variableDescription,
			type: queryType,
			queryValue: variableQueryValue,
			customValue: variableCustomValue,
			textboxValue: variableTextboxValue,
			multiSelect: variableMultiSelect,
			showALLOption: variableShowALLOption,
			sort: variableSortType,
			...(queryType === 'TEXTBOX' && {
				selectedValue: (variableData.selectedValue ||
					variableTextboxValue) as never,
			}),
			modificationUUID: v4(),
		};
		onSave(
			variableName,
			newVariableData,
			(variableViewMode === 'EDIT' && variableName !== variableData.name
				? variableData.name
				: '') as string,
		);
		onCancel();
	};

	// Fetches the preview values for the SQL variable query
	const handleQueryResult = async (): Promise<void> => {
		setPreviewLoading(true);
		setErrorPreview(null);
		try {
			const variableQueryResponse = await query({
				query: variableQueryValue,
				variables: variablePropsToPayloadVariables(existingVariables),
			});
			setPreviewLoading(false);
			if (variableQueryResponse.error) {
				let message = variableQueryResponse.error;
				if (variableQueryResponse.error.includes('Syntax error:')) {
					message =
						'Please make sure query is valid and dependent variables are selected';
				}
				setErrorPreview(message);
				return;
			}
			if (variableQueryResponse.payload?.variableValues)
				setPreviewValues(
					sortValues(
						variableQueryResponse.payload?.variableValues || [],
						variableSortType,
					) as never,
				);
		} catch (e) {
			console.error(e);
		}
	};
	return (
		<Col>
			{/* <Typography.Title level={3}>Add Variable</Typography.Title> */}
			<VariableItemRow>
				<LabelContainer>
					<Typography>Name</Typography>
				</LabelContainer>
				<div>
					<Input
						placeholder="Unique name of the variable"
						style={{ width: 400 }}
						value={variableName}
						onChange={(e): void => {
							setVariableName(e.target.value);
							setErrorName(
								!validateName(e.target.value) && e.target.value !== variableData.name,
							);
						}}
					/>
					<div>
						<Typography.Text type="warning">
							{errorName ? 'Variable name already exists' : ''}
						</Typography.Text>
					</div>
				</div>
			</VariableItemRow>
			<VariableItemRow>
				<LabelContainer>
					<Typography>Description</Typography>
				</LabelContainer>

				<Input.TextArea
					value={variableDescription}
					placeholder="Write description of the variable"
					style={{ width: 400 }}
					onChange={(e): void => setVariableDescription(e.target.value)}
				/>
			</VariableItemRow>
			<VariableItemRow>
				<LabelContainer>
					<Typography>Type</Typography>
				</LabelContainer>

				<Select
					defaultActiveFirstOption
					style={{ width: 400 }}
					onChange={(e: TVariableQueryType): void => {
						setQueryType(e);
					}}
					value={queryType}
				>
					<Option value={VariableQueryTypeArr[0]}>Query</Option>
					<Option value={VariableQueryTypeArr[1]}>Textbox</Option>
					<Option value={VariableQueryTypeArr[2]}>Custom</Option>
				</Select>
			</VariableItemRow>
			<Typography.Title
				level={5}
				style={{ marginTop: '1rem', marginBottom: '1rem' }}
			>
				Options
			</Typography.Title>
			{queryType === 'QUERY' && (
				<VariableItemRow>
					<LabelContainer>
						<Typography>Query</Typography>
					</LabelContainer>

					<div style={{ flex: 1, position: 'relative' }}>
						<Editor
							language="sql"
							value={variableQueryValue}
							onChange={(e): void => setVariableQueryValue(e)}
							height="300px"
						/>
						<Button
							type="primary"
							onClick={handleQueryResult}
							style={{
								position: 'absolute',
								bottom: 0,
							}}
							loading={previewLoading}
						>
							Test Run Query
						</Button>
					</div>
				</VariableItemRow>
			)}
			{queryType === 'CUSTOM' && (
				<VariableItemRow>
					<LabelContainer>
						<Typography>Values separated by comma</Typography>
					</LabelContainer>
					<Input.TextArea
						value={variableCustomValue}
						placeholder="1, 10, mykey, mykey:myvalue"
						style={{ width: 400 }}
						onChange={(e): void => {
							setVariableCustomValue(e.target.value);
							setPreviewValues(
								sortValues(
									commaValuesParser(e.target.value),
									variableSortType,
								) as never,
							);
						}}
					/>
				</VariableItemRow>
			)}
			{queryType === 'TEXTBOX' && (
				<VariableItemRow>
					<LabelContainer>
						<Typography>Default Value</Typography>
					</LabelContainer>
					<Input
						value={variableTextboxValue}
						onChange={(e): void => {
							setVariableTextboxValue(e.target.value);
						}}
						placeholder="Default value if any"
						style={{ width: 400 }}
					/>
				</VariableItemRow>
			)}
			{(queryType === 'QUERY' || queryType === 'CUSTOM') && (
				<>
					<VariableItemRow>
						<LabelContainer>
							<Typography>Preview of Values</Typography>
						</LabelContainer>
						<div style={{ flex: 1 }}>
							{errorPreview ? (
								<Typography style={{ color: orange[5] }}>{errorPreview}</Typography>
							) : (
								map(previewValues, (value, idx) => (
									<Tag key={`${value}${idx}`}>{value.toString()}</Tag>
								))
							)}
						</div>
					</VariableItemRow>
					<VariableItemRow>
						<LabelContainer>
							<Typography>Sort</Typography>
						</LabelContainer>

						<Select
							defaultActiveFirstOption
							style={{ width: 400 }}
							defaultValue={VariableSortTypeArr[0]}
							value={variableSortType}
							onChange={(value: TSortVariableValuesType): void =>
								setVariableSortType(value)
							}
						>
							<Option value={VariableSortTypeArr[0]}>Disabled</Option>
							<Option value={VariableSortTypeArr[1]}>Ascending</Option>
							<Option value={VariableSortTypeArr[2]}>Descending</Option>
						</Select>
					</VariableItemRow>
					<VariableItemRow>
						<LabelContainer>
							<Typography>Enable multiple values to be checked</Typography>
						</LabelContainer>
						<Switch
							checked={variableMultiSelect}
							onChange={(e): void => {
								setVariableMultiSelect(e);
								if (!e) {
									setVariableShowALLOption(false);
								}
							}}
						/>
					</VariableItemRow>
					{variableMultiSelect && (
						<VariableItemRow>
							<LabelContainer>
								<Typography>Include an option for ALL values</Typography>
							</LabelContainer>
							<Switch
								checked={variableShowALLOption}
								onChange={(e): void => setVariableShowALLOption(e)}
							/>
						</VariableItemRow>
					)}
				</>
			)}
			<Divider />
			<VariableItemRow>
				<Button type="primary" onClick={handleSave} disabled={errorName}>
					Save
				</Button>
				<Button type="dashed" onClick={onCancel}>
					Cancel
				</Button>
			</VariableItemRow>
		</Col>
	);
}

export default VariableItem;
