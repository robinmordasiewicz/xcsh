/**
 * Fix for intermittent sidebar visibility with navigation.instant
 *
 * Material for MkDocs' instant loading can cause race conditions where
 * the data-md-page attribute doesn't update correctly during navigation.
 * This causes the home page CSS rule to incorrectly hide the sidebar.
 *
 * This script ensures the data-md-page attribute is correctly set based on URL.
 */

// Only set up once - check if already initialized
if (!window.__sidebarFixInitialized) {
  window.__sidebarFixInitialized = true;

  function isHomePage(path) {
    return path.endsWith('/vesctl/') || path.endsWith('/vesctl/index.html') || path === '/';
  }

  function fixSidebar() {
    var path = window.location.pathname;
    var currentAttr = document.documentElement.getAttribute('data-md-page');
    var shouldBeHome = isHomePage(path);
    var isCurrentlyHome = currentAttr === 'home';

    // Only update if needed
    if (shouldBeHome && !isCurrentlyHome) {
      document.documentElement.setAttribute('data-md-page', 'home');
    } else if (!shouldBeHome && isCurrentlyHome) {
      document.documentElement.removeAttribute('data-md-page');
    }
  }

  // Run immediately
  fixSidebar();

  // Subscribe to Material for MkDocs navigation events
  // This is the key - document$ fires after each instant navigation
  if (typeof document$ !== 'undefined') {
    document$.subscribe(fixSidebar);
  }

  // Also handle browser back/forward
  window.addEventListener('popstate', fixSidebar);
}
