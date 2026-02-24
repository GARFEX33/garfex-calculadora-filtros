import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

/**
 * Combina clases CSS con soporte para condicionales y deduplicación Tailwind.
 * Usar siempre en lugar de concatenar strings de clases manualmente.
 *
 * @example
 * cn('px-4 py-2', isActive && 'bg-primary', className)
 */
export function cn(...inputs: ClassValue[]): string {
	return twMerge(clsx(inputs));
}
