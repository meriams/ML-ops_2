import {
	BooleanFormats,
	DataFormats,
	DataRateFormats,
	MiscellaneousFormats,
	ThroughputFormats,
	TimeFormats,
} from 'container/NewWidget/RightContainer/types';

import {
	dataFormatConfig,
	dataRateUnitsConfig,
	miscUnitsConfig,
	throughputConfig,
	timeUnitsConfig,
} from './config';

export function covertIntoDataFormats({
	value,
	sourceUnit,
	targetUnit,
}: IUnit): number {
	if (sourceUnit === undefined || targetUnit === undefined) {
		return value;
	}

	if (Object.values(BooleanFormats).includes(sourceUnit as BooleanFormats)) {
		return 1;
	}

	const sourceMultiplier =
		dataFormatConfig[sourceUnit as DataFormats] ||
		timeUnitsConfig[sourceUnit as TimeFormats] ||
		dataRateUnitsConfig[sourceUnit as DataRateFormats] ||
		miscUnitsConfig[sourceUnit as MiscellaneousFormats] ||
		throughputConfig[sourceUnit as ThroughputFormats];

	const targetDivider =
		dataFormatConfig[targetUnit as DataFormats] ||
		timeUnitsConfig[targetUnit as TimeFormats] ||
		dataRateUnitsConfig[targetUnit as DataRateFormats] ||
		miscUnitsConfig[targetUnit as MiscellaneousFormats] ||
		throughputConfig[sourceUnit as ThroughputFormats];

	const intermediateValue = value * sourceMultiplier;

	const roundedValue = Math.round(intermediateValue * 1000000) / 1000000;

	const result = roundedValue / targetDivider;

	return Number.isNaN(result) ? 0 : result;
}

interface IUnit {
	value: number;
	sourceUnit?: string;
	targetUnit?: string;
}
