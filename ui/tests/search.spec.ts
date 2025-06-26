import { test, expect } from '@playwright/test';

test('home shows search form', async ({ page }) => {
  await page.goto('/');
  await expect(page.getByPlaceholder('Vector e.g. 0.1,0.2')).toBeVisible();
});
