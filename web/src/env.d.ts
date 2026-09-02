/// <reference types="@rsbuild/core/types" />

declare namespace NodeJS {
	interface ProcessEnv {
		readonly NODE_ENV: 'development' | 'production' | 'test';
		/** 'true' enables the Tasks/Agenda surface. See src/config/features.ts. */
		readonly MINDSCAPE_FEATURE_TASKS?: string;
		/** 'true' enables the API Keys surface. See src/config/features.ts. */
		readonly MINDSCAPE_FEATURE_APIKEYS?: string;
		/** 'true' enables admin pages in a production build. */
		readonly MINDSCAPE_ADMIN?: string;
	}
}
