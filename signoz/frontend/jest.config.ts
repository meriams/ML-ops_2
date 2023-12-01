import type { Config } from '@jest/types';

const config: Config.InitialOptions = {
	clearMocks: true,
	coverageDirectory: 'coverage',
	coverageReporters: ['text', 'cobertura', 'html', 'json-summary'],
	moduleFileExtensions: ['ts', 'tsx', 'js', 'json'],
	modulePathIgnorePatterns: ['dist'],
	moduleNameMapper: {
		'\\.(css|less|scss)$': '<rootDir>/__mocks__/cssMock.ts',
	},
	globals: {
		extensionsToTreatAsEsm: ['.ts'],
		'ts-jest': {
			useESM: true,
		},
	},
	testMatch: ['<rootDir>/src/**/*?(*.)(test).(ts|js)?(x)'],
	preset: 'ts-jest/presets/js-with-ts-esm',
	transform: {
		'^.+\\.(ts|tsx)?$': 'ts-jest',
		'^.+\\.(js|jsx)$': 'babel-jest',
	},
	transformIgnorePatterns: [
		'node_modules/(?!(lodash-es|react-dnd|core-dnd|@react-dnd|dnd-core|react-dnd-html5-backend)/)',
	],
	setupFilesAfterEnv: ['<rootDir>jest.setup.ts'],
	testPathIgnorePatterns: ['/node_modules/', '/public/'],
	moduleDirectories: ['node_modules', 'src'],
	testEnvironment: 'jest-environment-jsdom',
	testEnvironmentOptions: {
		'jest-playwright': {
			browsers: ['chromium', 'firefox', 'webkit'],
		},
	},
};

export default config;
