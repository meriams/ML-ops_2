import { blue } from '@ant-design/colors';
import { Col, Row, Space } from 'antd';
import styled from 'styled-components';
import { getActiveLogBackground, getDefaultLogBackground } from 'utils/logs';

import { RawLogContentProps } from './types';

export const RawLogViewContainer = styled(Row)<{
	$isDarkMode: boolean;
	$isReadOnly?: boolean;
	$isActiveLog?: boolean;
}>`
	position: relative;
	width: 100%;
	font-weight: 700;
	font-size: 0.625rem;
	line-height: 1.25rem;

	transition: background-color 0.2s ease-in;

	${({ $isActiveLog }): string => getActiveLogBackground($isActiveLog)}

	${({ $isReadOnly, $isDarkMode, $isActiveLog }): string =>
		$isActiveLog
			? getActiveLogBackground()
			: getDefaultLogBackground($isReadOnly, $isDarkMode)}
`;

export const ExpandIconWrapper = styled(Col)`
	color: ${blue[6]};
	padding: 0.25rem 0.375rem;
	cursor: pointer;
	font-size: 12px;
`;

export const RawLogContent = styled.div<RawLogContentProps>`
	margin-bottom: 0;
	font-family: Fira Code, monospace;
	font-weight: 300;

	${({ $isTextOverflowEllipsisDisabled, linesPerRow }): string =>
		$isTextOverflowEllipsisDisabled
			? 'white-space: nowrap'
			: `overflow: hidden;
		text-overflow: ellipsis; 
		display: -webkit-box;
		-webkit-line-clamp: ${linesPerRow};
		line-clamp: ${linesPerRow}; 
		-webkit-box-orient: vertical;`};

	font-size: 1rem;
	line-height: 2rem;

	cursor: ${({ $isActiveLog, $isReadOnly }): string =>
		$isActiveLog || $isReadOnly ? 'initial' : 'pointer'};

	${({ $isActiveLog, $isReadOnly }): string =>
		$isReadOnly && $isActiveLog ? 'padding: 0 1.5rem;' : ''}
`;

export const ActionButtonsWrapper = styled(Space)`
	position: absolute;
	transform: translate(-50%, -50%);
	top: 50%;
	right: 0;
	cursor: pointer;
`;
