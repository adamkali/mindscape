import type { ChainIdentifier, RspackChain } from '@rsbuild/core';
import type { BabelConfigUtils, BabelLoaderOptions, BabelTransformOptions, PluginBabelOptions } from './types.js';
export declare const BABEL_JS_RULE = "babel-js";
export declare const castArray: <T>(arr?: T | T[]) => T[];
export declare const getBabelUtils: (config: BabelTransformOptions) => BabelConfigUtils;
export declare const applyUserBabelConfig: (defaultOptions: BabelLoaderOptions, userBabelConfig?: PluginBabelOptions["babelLoaderOptions"], extraBabelUtils?: Partial<BabelConfigUtils>) => BabelLoaderOptions;
export declare const modifyBabelLoaderOptions: ({ chain, CHAIN_ID, modifier, }: {
    chain: RspackChain;
    CHAIN_ID: ChainIdentifier;
    modifier: (config: BabelTransformOptions) => BabelTransformOptions;
}) => void;
