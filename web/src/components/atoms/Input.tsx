import {type ComponentProps, type JSX } from 'solid-js';

/**
 * Props for the Input component
 */
interface InputProps extends Omit<ComponentProps<'input'>, 'onInput'> {
  /** Visual variant of the input */
  variant?: 'primary' | 'secondary' | 'tertiary' | 'danger';
  /** Title/label to display above the input */
  title?: string;
  /** Current input value */
  value?: string;
  /** Callback function called when input value changes */
  onValueChange?: (value: string) => void;
  /** Whether the input has validation errors */
  error?: boolean;
  /** Error message to display below the input */
  errorMessage?: string;
  /** Icon or other elements to display alongside the input (typically from nerdfonts.com) */
  children?: JSX.Element;
}

/**
 * Input component with variant support and icon integration
 * 
 * Features:
 * - Multiple visual variants (primary, secondary, tertiary, danger)
  * - Title/label support
  * - Error state handling with validation messages
 * - Icon support from nerdfonts.com or other icon libraries
 * - Callback-based value management
 * - Full input element prop support
 * 
 * @example
 * ```tsx
 * <Input 
 *   variant="primary"
  *   title="Search"
 *   value={inputValue}
 *   onValueChange={setInputValue}
 *   placeholder="Enter text..."
  *   error={hasError}
  *   errorMessage="This field is required"
 * >
 *   <i class="nf-fa-search" />
 * </Input>
 * ```
 */
export default function Input(props: InputProps): JSX.Element {
  const { variant = 'primary', title, children, onValueChange, error, errorMessage, ...inputProps } = props;

  const handleInput = (e: InputEvent) => {
    const target = e.target as HTMLInputElement;
    onValueChange?.(target.value);
  };

  switch (variant) {
    case 'primary':
      return __primary_input(inputProps, children, handleInput, title, error, errorMessage);
    case 'secondary':
      return __secondary_input(inputProps, children, handleInput, title, error, errorMessage);
    case 'tertiary':
      return __tertiary_input(inputProps, children, handleInput, title, error, errorMessage);
    case 'danger':
      return __danger_input(inputProps, children, handleInput, title, error, errorMessage);
  }
}

function __primary_input(
  props: Omit<ComponentProps<'input'>, 'onInput'>,
  children?: JSX.Element,
  onInput?: (e: InputEvent) => void,
  title?: string,
  error?: boolean,
  errorMessage?: string
): JSX.Element {
  return (
    <div class="w-full">
      {title && (
        <label class={`block text-sm font-medium mb-2 ${error ? 'text-red-600' : 'text-primary'}`}>
          {title}
        </label>
      )}
      <div class={`flex items-center rounded-lg shadow-sm transition-all duration-200 ${
        error 
          ? 'bg-red-50 hover:bg-red-100 border-2 border-red-300 focus-within:border-red-500 focus-within:ring-2 focus-within:ring-red-200' 
          : 'bg-blue-600 hover:bg-blue-700 border-2 border-transparent focus-within:border-blue-300 focus-within:ring-2 focus-within:ring-blue-200'
      }`}>
        {children && (
          <div class={`pl-3 ${error ? 'text-red-500' : 'text-white/70'}`}>
            {children}
          </div>
        )}
        <input 
          {...props}
          onInput={onInput}
          class={`flex-1 bg-transparent py-3 px-4 font-medium focus:outline-none rounded-lg ${
            error 
              ? 'text-red-900 placeholder-red-400' 
              : 'text-white placeholder-white/50 bg-primary'
          }`}
        />
      </div>
      {error && errorMessage && (
        <p class="mt-1 text-sm text-red-600">{errorMessage}</p>
      )}
    </div>
  );
}

function __secondary_input(
  props: Omit<ComponentProps<'input'>, 'onInput'>,
  children?: JSX.Element,
  onInput?: (e: InputEvent) => void,
  title?: string,
  error?: boolean,
  errorMessage?: string
): JSX.Element {
  return (
    <div class="w-full">
      {title && (
        <label class={`block text-sm font-medium mb-2 ${error ? 'text-red-600' : 'text-slate-700 dark:text-slate-300'}`}>
          {title}
        </label>
      )}
      <div class={`flex items-center rounded-lg shadow-sm transition-all duration-200 ${
        error 
          ? 'bg-red-50 hover:bg-red-100 border-2 border-red-300 focus-within:border-red-500 focus-within:ring-2 focus-within:ring-red-200' 
          : 'bg-slate-100 hover:bg-slate-200 border-2 border-transparent focus-within:border-slate-400 focus-within:ring-2 focus-within:ring-slate-200'
      }`}>
        {children && (
          <div class={`pl-3 ${error ? 'text-red-500' : 'text-slate-600'}`}>
            {children}
          </div>
        )}
        <input 
          {...props}
          onInput={onInput}
          class={`flex-1 bg-transparent py-3 px-4 font-medium focus:outline-none rounded-lg ${
            error 
              ? 'text-red-900 placeholder-red-400' 
              : 'text-slate-900 placeholder-slate-500'
          }`}
        />
      </div>
      {error && errorMessage && (
        <p class="mt-1 text-sm text-red-600">{errorMessage}</p>
      )}
    </div>
  );
}

function __tertiary_input(
  props: Omit<ComponentProps<'input'>, 'onInput'>,
  children?: JSX.Element,
  onInput?: (e: InputEvent) => void,
  title?: string,
  error?: boolean,
  errorMessage?: string
): JSX.Element {
  return (
    <div class="w-full">
      {title && (
        <label class={`block text-sm font-medium mb-2 ${error ? 'text-red-600' : 'text-slate-700 dark:text-slate-300'}`}>
          {title}
        </label>
      )}
      <div class={`flex items-center rounded-lg transition-all duration-200 ${
        error 
          ? 'bg-red-50 hover:bg-red-100 border-2 border-red-300 focus-within:border-red-500 focus-within:ring-2 focus-within:ring-red-200' 
          : 'bg-white hover:bg-gray-50 border-2 border-gray-300 focus-within:border-gray-500 focus-within:ring-2 focus-within:ring-gray-200'
      }`}>
        {children && (
          <div class={`pl-3 ${error ? 'text-red-500' : 'text-gray-500'}`}>
            {children}
          </div>
        )}
        <input 
          {...props}
          onInput={onInput}
          class={`flex-1 bg-transparent py-3 px-4 font-medium focus:outline-none rounded-lg ${
            error 
              ? 'text-red-900 placeholder-red-400' 
              : 'text-gray-900 placeholder-gray-500'
          }`}
        />
      </div>
      {error && errorMessage && (
        <p class="mt-1 text-sm text-red-600">{errorMessage}</p>
      )}
    </div>
  );
}

function __danger_input(
  props: Omit<ComponentProps<'input'>, 'onInput'>,
  children?: JSX.Element,
  onInput?: (e: InputEvent) => void,
  title?: string,
  error?: boolean,
  errorMessage?: string
): JSX.Element {
  return (
    <div class="w-full">
      {title && (
        <label class="block text-sm font-medium mb-2 text-red-600">
          {title}
        </label>
      )}
      <div class={`flex items-center rounded-lg shadow-sm transition-all duration-200 ${
        error 
          ? 'bg-red-100 hover:bg-red-200 border-2 border-red-400 focus-within:border-red-600 focus-within:ring-2 focus-within:ring-red-300' 
          : 'bg-red-600 hover:bg-red-700 border-2 border-transparent focus-within:border-red-300 focus-within:ring-2 focus-within:ring-red-200'
      }`}>
        {children && (
          <div class={`pl-3 ${error ? 'text-red-600' : 'text-white/70'}`}>
            {children}
          </div>
        )}
        <input 
          {...props}
          onInput={onInput}
          class={`flex-1 bg-transparent py-3 px-4 font-medium focus:outline-none rounded-lg ${
            error 
              ? 'text-red-900 placeholder-red-500' 
              : 'text-white placeholder-white/50'
          }`}
        />
      </div>
      {error && errorMessage && (
        <p class="mt-1 text-sm text-red-600">{errorMessage}</p>
      )}
    </div>
  );
}

