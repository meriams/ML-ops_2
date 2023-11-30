import { DeleteOutlined, ExclamationCircleOutlined } from '@ant-design/icons';
import { Modal, Tooltip, Typography } from 'antd';
import { REACT_QUERY_KEY } from 'constants/reactQueryKeys';
import { useDeleteDashboard } from 'hooks/dashboard/useDeleteDashboard';
import { useNotifications } from 'hooks/useNotifications';
import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useQueryClient } from 'react-query';
import { useSelector } from 'react-redux';
import { AppState } from 'store/reducers';
import AppReducer from 'types/reducer/app';
import { USER_ROLES } from 'types/roles';

import { Data } from '..';
import { TableLinkText } from './styles';

interface DeleteButtonProps {
	createdBy: string;
	name: string;
	id: string;
	isLocked: boolean;
}

function DeleteButton({
	createdBy,
	name,
	id,
	isLocked,
}: DeleteButtonProps): JSX.Element {
	const [modal, contextHolder] = Modal.useModal();
	const { role, user } = useSelector<AppState, AppReducer>((state) => state.app);
	const isAuthor = user?.email === createdBy;

	const queryClient = useQueryClient();

	const { notifications } = useNotifications();

	const { t } = useTranslation(['dashboard']);

	const deleteDashboardMutation = useDeleteDashboard(id);

	const openConfirmationDialog = useCallback((): void => {
		modal.confirm({
			title: (
				<Typography.Title level={5}>
					Are you sure you want to delete the
					<span style={{ color: '#e42b35', fontWeight: 500 }}> {name} </span>
					dashboard?
				</Typography.Title>
			),
			icon: <ExclamationCircleOutlined style={{ color: '#e42b35' }} />,
			onOk() {
				deleteDashboardMutation.mutateAsync(undefined, {
					onSuccess: () => {
						notifications.success({
							message: t('dashboard:delete_dashboard_success', {
								name,
							}),
						});
						queryClient.invalidateQueries([REACT_QUERY_KEY.GET_ALL_DASHBOARDS]);
					},
				});
			},
			okText: 'Delete',
			okButtonProps: { danger: true },
			centered: true,
		});
	}, [modal, name, deleteDashboardMutation, notifications, t, queryClient]);

	const getDeleteTooltipContent = (): string => {
		if (isLocked) {
			if (role === USER_ROLES.ADMIN || isAuthor) {
				return t('dashboard:locked_dashboard_delete_tooltip_admin_author');
			}

			return t('dashboard:locked_dashboard_delete_tooltip_editor');
		}

		return '';
	};

	return (
		<>
			<Tooltip placement="left" title={getDeleteTooltipContent()}>
				<TableLinkText
					type="danger"
					onClick={(): void => {
						if (!isLocked) {
							openConfirmationDialog();
						}
					}}
					disabled={isLocked}
				>
					<DeleteOutlined /> Delete
				</TableLinkText>
			</Tooltip>

			{contextHolder}
		</>
	);
}

// This is to avoid the type collision
function Wrapper(props: Data): JSX.Element {
	const {
		createdAt,
		description,
		id,
		key,
		lastUpdatedTime,
		name,
		tags,
		createdBy,
		lastUpdatedBy,
		isLocked,
	} = props;

	return (
		<DeleteButton
			{...{
				createdAt,
				description,
				id,
				key,
				lastUpdatedTime,
				name,
				tags,
				createdBy,
				lastUpdatedBy,
				isLocked,
			}}
		/>
	);
}

export default Wrapper;
