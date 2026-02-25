import type { NormalizedEnvironmentConfig, RsbuildContext, RsbuildPlugin } from '@rsbuild/core';
import type { BabelLoaderOptions, PluginBabelOptions } from './types.js';
export declare const PLUGIN_BABEL_NAME = "rsbuild:babel";
/**
 * The `@babel/preset-typescript` default options.
 */
export declare const DEFAULT_BABEL_PRESET_TYPESCRIPT_OPTIONS: {
    allowNamespaces: boolean;
    allExtensions: boolean;
    allowDeclareFields: boolean;
    optimizeConstEnums: boolean;
    isTSX: boolean;
};
export declare const getDefaultBabelOptions: (config: NormalizedEnvironmentConfig, context: RsbuildContext) => BabelLoaderOptions;
export declare const pluginBabel: (options?: PluginBabelOptions) => RsbuildPlugin;
