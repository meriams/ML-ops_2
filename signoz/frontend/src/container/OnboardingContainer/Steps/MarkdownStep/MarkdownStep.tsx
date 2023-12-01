/* eslint-disable @typescript-eslint/ban-ts-comment */
import { MarkdownRenderer } from 'components/MarkdownRenderer/MarkdownRenderer';
import { ApmDocFilePaths } from 'container/OnboardingContainer/constants/apmDocFilePaths';
import { InfraMonitoringDocFilePaths } from 'container/OnboardingContainer/constants/infraMonitoringDocFilePaths';
import { LogsManagementDocFilePaths } from 'container/OnboardingContainer/constants/logsManagementDocFilePaths';
import { useOnboardingContext } from 'container/OnboardingContainer/context/OnboardingContext';
import { ModulesMap } from 'container/OnboardingContainer/OnboardingContainer';
import useAnalytics from 'hooks/analytics/useAnalytics';
import { useEffect, useState } from 'react';

export interface IngestionInfoProps {
	SIGNOZ_INGESTION_KEY?: string;
	REGION?: string;
}

export default function MarkdownStep(): JSX.Element {
	const {
		activeStep,
		ingestionData,
		serviceName,
		selectedDataSource,
		selectedModule,
		selectedEnvironment,
		selectedFramework,
		selectedMethod,
	} = useOnboardingContext();

	const { trackEvent } = useAnalytics();

	const [markdownContent, setMarkdownContent] = useState('');

	const { step } = activeStep;

	const getFilePath = (): any => {
		let path = `${selectedModule?.id}_${selectedDataSource?.id}`;

		if (selectedFramework) {
			path += `_${selectedFramework}`;
		}

		if (selectedEnvironment) {
			path += `_${selectedEnvironment}`;
		}

		if (
			selectedModule?.id === ModulesMap.APM &&
			selectedDataSource?.id !== 'kubernetes' &&
			selectedMethod
		) {
			path += `_${selectedMethod}`;
		}

		path += `_${step?.id}`;

		return path;
	};

	useEffect(() => {
		const path = getFilePath();

		let docFilePaths;

		if (selectedModule?.id === ModulesMap.APM) {
			docFilePaths = ApmDocFilePaths;
		} else if (selectedModule?.id === ModulesMap.LogsManagement) {
			docFilePaths = LogsManagementDocFilePaths;
		} else if (selectedModule?.id === ModulesMap.InfrastructureMonitoring) {
			docFilePaths = InfraMonitoringDocFilePaths;
		}

		// @ts-ignore
		if (docFilePaths && docFilePaths[path]) {
			// @ts-ignore
			setMarkdownContent(docFilePaths[path]);
		}

		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [step]);

	const variables = {
		MYAPP: serviceName || '<service-name>',
		SIGNOZ_INGESTION_KEY:
			ingestionData?.SIGNOZ_INGESTION_KEY || '<SIGNOZ_INGESTION_KEY>',
		REGION: ingestionData?.REGION || 'region',
	};

	useEffect(() => {
		trackEvent(
			`Onboarding: ${activeStep?.module?.id}: ${selectedDataSource?.name}: ${activeStep?.step?.title}`,
			{
				dataSource: selectedDataSource,
				framework: selectedFramework,
				environment: selectedEnvironment,
				module: {
					name: activeStep?.module?.title,
					id: activeStep?.module?.id,
				},
				step: {
					name: activeStep?.step?.title,
					id: activeStep?.step?.id,
				},
			},
		);
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [step]);

	return (
		<div className="markdown-container">
			<MarkdownRenderer markdownContent={markdownContent} variables={variables} />
		</div>
	);
}
