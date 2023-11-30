/* eslint-disable react/display-name */
import { PlusOutlined } from '@ant-design/icons';
import { Typography } from 'antd';
import { ColumnsType } from 'antd/lib/table';
import saveAlertApi from 'api/alerts/save';
import DropDown from 'components/DropDown/DropDown';
import {
	DynamicColumnsKey,
	TableDataSource,
} from 'components/ResizeTable/contants';
import DynamicColumnTable from 'components/ResizeTable/DynamicColumnTable';
import DateComponent from 'components/ResizeTable/TableComponent/DateComponent';
import LabelColumn from 'components/TableRenderer/LabelColumn';
import TextToolTip from 'components/TextToolTip';
import { QueryParams } from 'constants/query';
import ROUTES from 'constants/routes';
import useComponentPermission from 'hooks/useComponentPermission';
import useInterval from 'hooks/useInterval';
import { useNotifications } from 'hooks/useNotifications';
import history from 'lib/history';
import { mapQueryDataFromApi } from 'lib/newQueryBuilder/queryBuilderMappers/mapQueryDataFromApi';
import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { UseQueryResult } from 'react-query';
import { useSelector } from 'react-redux';
import { AppState } from 'store/reducers';
import { ErrorResponse, SuccessResponse } from 'types/api';
import { GettableAlert } from 'types/api/alerts/get';
import AppReducer from 'types/reducer/app';

import DeleteAlert from './DeleteAlert';
import { Button, ButtonContainer, ColumnButton } from './styles';
import Status from './TableComponents/Status';
import ToggleAlertState from './ToggleAlertState';

function ListAlert({ allAlertRules, refetch }: ListAlertProps): JSX.Element {
	const [data, setData] = useState<GettableAlert[]>(allAlertRules || []);
	const { t } = useTranslation('common');
	const { role, featureResponse } = useSelector<AppState, AppReducer>(
		(state) => state.app,
	);
	const [addNewAlert, action] = useComponentPermission(
		['add_new_alert', 'action'],
		role,
	);

	const { notifications: notificationsApi } = useNotifications();

	useInterval(() => {
		(async (): Promise<void> => {
			const { data: refetchData, status } = await refetch();
			if (status === 'success') {
				setData(refetchData?.payload || []);
			}
			if (status === 'error') {
				notificationsApi.error({
					message: t('something_went_wrong'),
				});
			}
		})();
	}, 30000);

	const handleError = useCallback((): void => {
		notificationsApi.error({
			message: t('something_went_wrong'),
		});
	}, [notificationsApi, t]);

	const onClickNewAlertHandler = useCallback(() => {
		featureResponse
			.refetch()
			.then(() => {
				history.push(ROUTES.ALERTS_NEW);
			})
			.catch(handleError);
	}, [featureResponse, handleError]);

	const onEditHandler = (record: GettableAlert) => (): void => {
		featureResponse
			.refetch()
			.then(() => {
				const compositeQuery = mapQueryDataFromApi(record.condition.compositeQuery);

				history.push(
					`${ROUTES.EDIT_ALERTS}?ruleId=${record.id.toString()}&${
						QueryParams.compositeQuery
					}=${encodeURIComponent(JSON.stringify(compositeQuery))}`,
				);
			})
			.catch(handleError);
	};

	const onCloneHandler = (
		originalAlert: GettableAlert,
	) => async (): Promise<void> => {
		const copyAlert = {
			...originalAlert,
			alert: originalAlert.alert.concat(' - Copy'),
		};
		const apiReq = { data: copyAlert };

		const response = await saveAlertApi(apiReq);

		if (response.statusCode === 200) {
			notificationsApi.success({
				message: 'Success',
				description: 'Alert cloned successfully',
			});

			const { data: refetchData, status } = await refetch();
			if (status === 'success' && refetchData.payload) {
				setData(refetchData.payload || []);
				setTimeout(() => {
					const clonedAlert = refetchData.payload[refetchData.payload.length - 1];
					history.push(`${ROUTES.EDIT_ALERTS}?ruleId=${clonedAlert.id}`);
				}, 2000);
			}
			if (status === 'error') {
				notificationsApi.error({
					message: t('something_went_wrong'),
				});
			}
		} else {
			notificationsApi.error({
				message: 'Error',
				description: response.error || t('something_went_wrong'),
			});
		}
	};

	const dynamicColumns: ColumnsType<GettableAlert> = [
		{
			title: 'Created At',
			dataIndex: 'createAt',
			width: 80,
			key: DynamicColumnsKey.CreatedAt,
			align: 'center',
			sorter: (a: GettableAlert, b: GettableAlert): number => {
				const prev = new Date(a.createAt).getTime();
				const next = new Date(b.createAt).getTime();

				return prev - next;
			},
			render: DateComponent,
		},
		{
			title: 'Created By',
			dataIndex: 'createBy',
			width: 80,
			key: DynamicColumnsKey.CreatedBy,
			align: 'center',
		},
		{
			title: 'Updated At',
			dataIndex: 'updateAt',
			width: 80,
			key: DynamicColumnsKey.UpdatedAt,
			align: 'center',
			sorter: (a: GettableAlert, b: GettableAlert): number => {
				const prev = new Date(a.updateAt).getTime();
				const next = new Date(b.updateAt).getTime();

				return prev - next;
			},
			render: DateComponent,
		},
		{
			title: 'Updated By',
			dataIndex: 'updateBy',
			width: 80,
			key: DynamicColumnsKey.UpdatedBy,
			align: 'center',
		},
	];

	const columns: ColumnsType<GettableAlert> = [
		{
			title: 'Status',
			dataIndex: 'state',
			width: 80,
			key: 'state',
			sorter: (a, b): number =>
				(b.state ? b.state.charCodeAt(0) : 1000) -
				(a.state ? a.state.charCodeAt(0) : 1000),
			render: (value): JSX.Element => <Status status={value} />,
		},
		{
			title: 'Alert Name',
			dataIndex: 'alert',
			width: 100,
			key: 'name',
			sorter: (alertA, alertB): number => {
				if (alertA.alert && alertB.alert) {
					return alertA.alert.localeCompare(alertB.alert);
				}
				return 0;
			},
			render: (value, record): JSX.Element => (
				<Typography.Link onClick={onEditHandler(record)}>{value}</Typography.Link>
			),
		},
		{
			title: 'Severity',
			dataIndex: 'labels',
			width: 80,
			key: 'severity',
			sorter: (a, b): number =>
				(a.labels ? a.labels.severity.length : 0) -
				(b.labels ? b.labels.severity.length : 0),
			render: (value): JSX.Element => {
				const objectKeys = Object.keys(value);
				const withSeverityKey = objectKeys.find((e) => e === 'severity') || '';
				const severityValue = value[withSeverityKey];

				return <Typography>{severityValue}</Typography>;
			},
		},
		{
			title: 'Labels',
			dataIndex: 'labels',
			key: 'tags',
			align: 'center',
			width: 100,
			render: (value): JSX.Element => {
				const objectKeys = Object.keys(value);
				const withOutSeverityKeys = objectKeys.filter((e) => e !== 'severity');

				if (withOutSeverityKeys.length === 0) {
					return <Typography>-</Typography>;
				}

				return (
					<LabelColumn labels={withOutSeverityKeys} value={value} color="magenta" />
				);
			},
		},
	];

	if (action) {
		columns.push({
			title: 'Action',
			dataIndex: 'id',
			key: 'action',
			width: 10,
			render: (id: GettableAlert['id'], record): JSX.Element => (
				<DropDown
					element={[
						<ToggleAlertState
							key="1"
							disabled={record.disabled}
							setData={setData}
							id={id}
						/>,
						<ColumnButton key="2" onClick={onEditHandler(record)} type="link">
							Edit
						</ColumnButton>,
						<ColumnButton key="3" onClick={onCloneHandler(record)} type="link">
							Clone
						</ColumnButton>,
						<DeleteAlert
							key="4"
							notifications={notificationsApi}
							setData={setData}
							id={id}
						/>,
					]}
				/>
			),
		});
	}

	return (
		<>
			<ButtonContainer>
				<TextToolTip
					{...{
						text: `More details on how to create alerts`,
						url: 'https://signoz.io/docs/userguide/alerts-management/',
					}}
				/>

				{addNewAlert && (
					<Button onClick={onClickNewAlertHandler} icon={<PlusOutlined />}>
						New Alert
					</Button>
				)}
			</ButtonContainer>
			<DynamicColumnTable
				tablesource={TableDataSource.Alert}
				columns={columns}
				rowKey="id"
				dataSource={data}
				dynamicColumns={dynamicColumns}
			/>
		</>
	);
}

interface ListAlertProps {
	allAlertRules: GettableAlert[];
	refetch: UseQueryResult<
		ErrorResponse | SuccessResponse<GettableAlert[]>
	>['refetch'];
}

export default ListAlert;
