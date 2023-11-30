import { fireEvent, render, screen } from '@testing-library/react';
import ROUTES from 'constants/routes';
import MockQueryClientProvider from 'providers/test/MockQueryClientProvider';
import { DataSource } from 'types/common/queryBuilder';

import { viewMockData } from '../__mock__/viewData';
import ExplorerCard from '../ExplorerCard';

jest.mock('react-router-dom', () => ({
	...jest.requireActual('react-router-dom'),
	useLocation: (): { pathname: string } => ({
		pathname: `${process.env.FRONTEND_API_ENDPOINT}/${ROUTES.TRACES_EXPLORER}/`,
	}),
}));

jest.mock('hooks/queryBuilder/useGetPanelTypesQueryParam', () => ({
	useGetPanelTypesQueryParam: jest.fn(() => 'mockedPanelType'),
}));

jest.mock('hooks/saveViews/useGetAllViews', () => ({
	useGetAllViews: jest.fn(() => ({
		data: { data: { data: viewMockData } },
		isLoading: false,
		error: null,
		isRefetching: false,
		refetch: jest.fn(),
	})),
}));

jest.mock('hooks/saveViews/useUpdateView', () => ({
	useUpdateView: jest.fn(() => ({
		mutateAsync: jest.fn(),
	})),
}));

jest.mock('hooks/saveViews/useDeleteView', () => ({
	useDeleteView: jest.fn(() => ({
		mutateAsync: jest.fn(),
	})),
}));

describe('ExplorerCard', () => {
	it('renders a card with a title and a description', () => {
		render(
			<MockQueryClientProvider>
				<ExplorerCard sourcepage={DataSource.TRACES}>child</ExplorerCard>
			</MockQueryClientProvider>,
		);
		expect(screen.getByText('Query Builder')).toBeInTheDocument();
	});

	it('renders a save view button', () => {
		render(
			<MockQueryClientProvider>
				<ExplorerCard sourcepage={DataSource.TRACES}>child</ExplorerCard>
			</MockQueryClientProvider>,
		);
		expect(screen.getByText('Save view')).toBeInTheDocument();
	});

	it('should see all the view listed in dropdown', async () => {
		const screen = render(
			<ExplorerCard sourcepage={DataSource.TRACES}>Mock Children</ExplorerCard>,
		);
		const selectPlaceholder = screen.getByText('Select a view');

		fireEvent.mouseDown(selectPlaceholder);
		const viewNameText = await screen.getAllByText('View 1');
		viewNameText.forEach((element) => {
			expect(element).toBeInTheDocument();
		});
	});
});
