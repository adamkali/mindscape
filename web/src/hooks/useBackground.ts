import { createMemo } from 'solid-js';
import { useBackground as useBackgroundContext } from '@/contexts/BackgroundContext';

export const useBackground = useBackgroundContext;

export const useBackgroundStyle = () => {
	const { currentBackground } = useBackground();
	
	return createMemo(() => ({
		'background-image': `url(${currentBackground()})`,
		'background-size': 'cover',
		'background-position': 'center center',
		'background-repeat': 'no-repeat',
		'background-attachment': 'fixed'
	}));
};

export const useBackgroundClass = () => {
	const { currentBackground } = useBackground();
	
	return createMemo(() => 
		currentBackground() ? 'min-h-screen bg-background' : 'min-h-screen bg-background bg-gray-900'
	);
};