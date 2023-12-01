import { Typography } from 'antd';
import { useIsDarkMode } from 'hooks/useDarkMode';
import { memo } from 'react';

// ** Types
import { FilterLabelProps } from './FilterLabel.interfaces';
// ** Styles
import { StyledLabel } from './FilterLabel.styled';

export const FilterLabel = memo(function FilterLabel({
	label,
}: FilterLabelProps): JSX.Element {
	const isDarkMode = useIsDarkMode();

	return (
		<StyledLabel isDarkMode={isDarkMode}>
			<Typography>{label}</Typography>
		</StyledLabel>
	);
});
