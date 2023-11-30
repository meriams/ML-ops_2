import { InputNumber } from 'antd';
import styled from 'styled-components';

export const MaxLinesFieldWrapper = styled.div`
	display: flex;
	justify-content: space-between;
	align-items: center;
`;

export const MaxLinesInput = styled(InputNumber)`
	max-width: 46px;
`;
