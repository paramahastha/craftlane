import { test, expect } from '@playwright/test';

test('dashboard has provider observation surfaces', async ({ page }) => {
  await page.goto('/');
  await expect(page.getByRole('heading', { name: 'Craftlane' })).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Provider coverage' })).toBeVisible();
});
