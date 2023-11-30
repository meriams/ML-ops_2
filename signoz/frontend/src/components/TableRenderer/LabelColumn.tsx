import './LabelColumn.styles.scss';

import { Popover, Tag } from 'antd';
import { popupContainer } from 'utils/selectPopupContainer';

import { LabelColumnProps } from './TableRenderer.types';
import TagWithToolTip from './TagWithToolTip';

function LabelColumn({ labels, value, color }: LabelColumnProps): JSX.Element {
	const newLabels = labels.length > 3 ? labels.slice(0, 3) : labels;
	const remainingLabels = labels.length > 3 ? labels.slice(3) : [];

	return (
		<div className="label-column">
			{newLabels.map(
				(label: string): JSX.Element => (
					<TagWithToolTip key={label} label={label} color={color} value={value} />
				),
			)}
			{remainingLabels.length > 0 && (
				<Popover
					getPopupContainer={popupContainer}
					placement="bottomRight"
					showArrow={false}
					content={
						<div>
							{labels.map(
								(label: string): JSX.Element => (
									<TagWithToolTip
										key={label}
										label={label}
										color={color}
										value={value}
									/>
								),
							)}
						</div>
					}
					trigger="hover"
				>
					<Tag className="label-column--tag" color={color}>
						+{remainingLabels.length}
					</Tag>
				</Popover>
			)}
		</div>
	);
}

LabelColumn.defaultProps = {
	value: {},
};

export default LabelColumn;
