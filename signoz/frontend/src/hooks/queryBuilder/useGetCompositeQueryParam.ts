import { QueryParams } from 'constants/query';
import useUrlQuery from 'hooks/useUrlQuery';
import { useMemo } from 'react';
import { Query } from 'types/api/queryBuilder/queryBuilderData';

export const useGetCompositeQueryParam = (): Query | null => {
	const urlQuery = useUrlQuery();

	return useMemo(() => {
		const compositeQuery = urlQuery.get(QueryParams.compositeQuery);
		let parsedCompositeQuery: Query | null = null;

		try {
			if (!compositeQuery) return null;

			parsedCompositeQuery = JSON.parse(decodeURIComponent(compositeQuery));
		} catch (e) {
			parsedCompositeQuery = null;
		}

		return parsedCompositeQuery;
	}, [urlQuery]);
};
