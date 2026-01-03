/**
 * Subscription API Client
 *
 * Client for interacting with F5 Distributed Cloud subscription, billing, and quota APIs
 */

import type { APIClient } from "../../api/client.js";
import type {
	UsagePlan,
	UsagePlanListResponse,
	AddonService,
	AddonServiceListResponse,
	AddonSubscription,
	AddonSubscriptionListResponse,
	AddonActivationStatus,
	QuotaLimitsResponse,
	QuotaUsageResponse,
	CurrentUsage,
	HourlyUsageDetails,
	MonthlyUsage,
	MonthlyUsageListResponse,
	PaymentMethod,
	PaymentMethodListResponse,
	Invoice,
	InvoiceListResponse,
	PlanTransitionRequest,
	PlanTransitionResponse,
	AddonSubscribeRequest,
	AddonUnsubscribeRequest,
	SubscriptionOverview,
} from "./types.js";

/**
 * Subscription API Client
 *
 * Provides methods for managing subscriptions, addons, quotas, usage, and billing
 */
export class SubscriptionClient {
	constructor(private apiClient: APIClient) {}

	// ============ Usage Plans ============

	/**
	 * Get current usage plan
	 */
	async getCurrentPlan(): Promise<UsagePlan> {
		const response = await this.apiClient.get<UsagePlan>(
			"/api/web/namespaces/system/usage_plans/current",
		);

		if (!response.ok) {
			throw new Error(
				`Failed to get current plan: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data;
	}

	/**
	 * List all available usage plans
	 */
	async listPlans(): Promise<UsagePlan[]> {
		const response = await this.apiClient.get<UsagePlanListResponse>(
			"/api/web/namespaces/system/usage_plans/custom_list",
		);

		if (!response.ok) {
			throw new Error(
				`Failed to list plans: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data.items ?? [];
	}

	/**
	 * Request plan transition
	 */
	async transitionPlan(
		request: PlanTransitionRequest,
	): Promise<PlanTransitionResponse> {
		const response = await this.apiClient.post<PlanTransitionResponse>(
			"/api/web/namespaces/system/billing/plan_transition",
			request as unknown as Record<string, unknown>,
		);

		if (!response.ok) {
			throw new Error(
				`Failed to transition plan: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data;
	}

	// ============ Addon Services ============

	/**
	 * List all addon services
	 */
	async listAddonServices(
		namespace: string = "system",
	): Promise<AddonService[]> {
		const response = await this.apiClient.get<AddonServiceListResponse>(
			`/api/web/namespaces/${namespace}/addon_services`,
		);

		if (!response.ok) {
			throw new Error(
				`Failed to list addon services: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data.items ?? [];
	}

	/**
	 * Get addon service details
	 */
	async getAddonService(
		name: string,
		namespace: string = "system",
	): Promise<AddonService> {
		const response = await this.apiClient.get<AddonService>(
			`/api/web/namespaces/${namespace}/addon_services/${name}`,
		);

		if (!response.ok) {
			throw new Error(
				`Failed to get addon service: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data;
	}

	/**
	 * List addon subscriptions
	 */
	async listAddonSubscriptions(
		namespace: string = "system",
	): Promise<AddonSubscription[]> {
		const response =
			await this.apiClient.get<AddonSubscriptionListResponse>(
				`/api/web/namespaces/${namespace}/addon_subscriptions`,
			);

		if (!response.ok) {
			throw new Error(
				`Failed to list addon subscriptions: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data.items ?? [];
	}

	/**
	 * Get addon activation status
	 */
	async getAddonActivationStatus(
		addonName: string,
	): Promise<AddonActivationStatus> {
		const response = await this.apiClient.get<AddonActivationStatus>(
			`/api/web/namespaces/system/addon_services/${addonName}/activation-status`,
		);

		if (!response.ok) {
			throw new Error(
				`Failed to get addon status: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data;
	}

	/**
	 * Get all addon activation statuses
	 */
	async getAllAddonActivationStatus(): Promise<AddonActivationStatus[]> {
		const response = await this.apiClient.get<{
			items?: AddonActivationStatus[];
		}>("/api/web/namespaces/system/addon_services/all-activation-status");

		if (!response.ok) {
			throw new Error(
				`Failed to get all addon statuses: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data.items ?? [];
	}

	/**
	 * Subscribe to an addon
	 */
	async subscribeToAddon(request: AddonSubscribeRequest): Promise<void> {
		const response = await this.apiClient.post<Record<string, unknown>>(
			"/api/web/namespaces/system/addon/subscribe",
			request as unknown as Record<string, unknown>,
		);

		if (!response.ok) {
			throw new Error(
				`Failed to subscribe to addon: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}
	}

	/**
	 * Unsubscribe from an addon
	 */
	async unsubscribeFromAddon(
		request: AddonUnsubscribeRequest,
	): Promise<void> {
		const response = await this.apiClient.post<Record<string, unknown>>(
			"/api/web/namespaces/system/addon/unsubscribe",
			request as unknown as Record<string, unknown>,
		);

		if (!response.ok) {
			throw new Error(
				`Failed to unsubscribe from addon: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}
	}

	// ============ Quotas ============

	/**
	 * Get quota limits
	 */
	async getQuotaLimits(
		namespace: string = "system",
	): Promise<QuotaLimitsResponse> {
		const response = await this.apiClient.get<QuotaLimitsResponse>(
			`/api/web/namespaces/${namespace}/quota/limits`,
		);

		if (!response.ok) {
			throw new Error(
				`Failed to get quota limits: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data;
	}

	/**
	 * Get quota usage
	 */
	async getQuotaUsage(
		namespace: string = "system",
	): Promise<QuotaUsageResponse> {
		const response = await this.apiClient.get<QuotaUsageResponse>(
			`/api/web/namespaces/${namespace}/quota/usage`,
		);

		if (!response.ok) {
			throw new Error(
				`Failed to get quota usage: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data;
	}

	// ============ Usage ============

	/**
	 * Get current usage
	 */
	async getCurrentUsage(namespace: string = "system"): Promise<CurrentUsage> {
		const response = await this.apiClient.get<CurrentUsage>(
			`/api/web/namespaces/${namespace}/current_usage`,
		);

		if (!response.ok) {
			throw new Error(
				`Failed to get current usage: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data;
	}

	/**
	 * Get hourly usage details
	 */
	async getHourlyUsage(
		namespace: string = "system",
		date?: string,
	): Promise<HourlyUsageDetails> {
		let url = `/api/web/namespaces/${namespace}/hourly_usage_details`;
		if (date) {
			url += `?date=${date}`;
		}

		const response = await this.apiClient.get<HourlyUsageDetails>(url);

		if (!response.ok) {
			throw new Error(
				`Failed to get hourly usage: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data;
	}

	/**
	 * Get monthly usage
	 */
	async getMonthlyUsage(
		namespace: string = "system",
	): Promise<MonthlyUsage[]> {
		const response = await this.apiClient.get<MonthlyUsageListResponse>(
			`/api/web/namespaces/${namespace}/monthly_usage`,
		);

		if (!response.ok) {
			throw new Error(
				`Failed to get monthly usage: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data.items ?? [];
	}

	/**
	 * Get detailed usage data
	 */
	async getUsageDetails(namespace: string = "system"): Promise<CurrentUsage> {
		const response = await this.apiClient.get<CurrentUsage>(
			`/api/web/namespaces/${namespace}/usage_details`,
		);

		if (!response.ok) {
			throw new Error(
				`Failed to get usage details: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data;
	}

	// ============ Billing / Payment Methods ============

	/**
	 * List payment methods
	 */
	async listPaymentMethods(
		namespace: string = "system",
	): Promise<PaymentMethod[]> {
		const response = await this.apiClient.get<PaymentMethodListResponse>(
			`/api/web/namespaces/${namespace}/billing/payment_methods`,
		);

		if (!response.ok) {
			throw new Error(
				`Failed to list payment methods: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data.items ?? [];
	}

	/**
	 * Get payment method details
	 */
	async getPaymentMethod(
		name: string,
		namespace: string = "system",
	): Promise<PaymentMethod> {
		const response = await this.apiClient.get<PaymentMethod>(
			`/api/web/namespaces/${namespace}/billing/payment_methods/${name}`,
		);

		if (!response.ok) {
			throw new Error(
				`Failed to get payment method: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data;
	}

	/**
	 * Set payment method as primary
	 */
	async setPaymentMethodPrimary(
		name: string,
		namespace: string = "system",
	): Promise<void> {
		const response = await this.apiClient.post<Record<string, unknown>>(
			`/api/web/namespaces/${namespace}/billing/payment_method/${name}/primary`,
			{},
		);

		if (!response.ok) {
			throw new Error(
				`Failed to set primary payment method: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}
	}

	/**
	 * Delete payment method
	 */
	async deletePaymentMethod(
		name: string,
		namespace: string = "system",
	): Promise<void> {
		const response = await this.apiClient.delete<Record<string, unknown>>(
			`/api/web/namespaces/${namespace}/billing/payment_methods/${name}`,
		);

		if (!response.ok) {
			throw new Error(
				`Failed to delete payment method: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}
	}

	// ============ Invoices ============

	/**
	 * List invoices
	 */
	async listInvoices(namespace: string = "system"): Promise<Invoice[]> {
		const response = await this.apiClient.get<InvoiceListResponse>(
			`/api/web/namespaces/${namespace}/usage/invoices/custom_list`,
		);

		if (!response.ok) {
			throw new Error(
				`Failed to list invoices: ${response.statusCode} - ${JSON.stringify(response.data)}`,
			);
		}

		return response.data.items ?? [];
	}

	/**
	 * Get invoice PDF URL path
	 * Note: Returns the relative path to download the invoice PDF
	 */
	getInvoicePdfPath(invoiceId: string, namespace: string = "system"): string {
		return `/api/web/namespaces/${namespace}/usage/invoice_pdf?invoice_id=${invoiceId}`;
	}

	// ============ Overview / Reports ============

	/**
	 * Get subscription overview (combines multiple API calls)
	 */
	async getOverview(): Promise<SubscriptionOverview> {
		// Fetch data in parallel for efficiency
		const [plan, addons, quotaUsage, currentUsage, paymentMethods] =
			await Promise.allSettled([
				this.getCurrentPlan(),
				this.listAddonServices(),
				this.getQuotaUsage(),
				this.getCurrentUsage(),
				this.listPaymentMethods(),
			]);

		const overview: SubscriptionOverview = {};

		// Plan
		if (plan.status === "fulfilled") {
			overview.plan = plan.value;
		}

		// Addon summary
		if (addons.status === "fulfilled") {
			const addonList = addons.value;
			const active = addonList.filter(
				(a) =>
					a.access_type === "SUBSCRIBED" ||
					a.access_type === "ENABLED",
			).length;
			overview.addon_summary = {
				active,
				available: addonList.filter((a) => a.status === "AVAILABLE")
					.length,
				total: addonList.length,
			};
		}

		// Quota summary
		if (quotaUsage.status === "fulfilled") {
			const usage = quotaUsage.value.usage ?? [];
			const critical = usage.filter(
				(q) => q.percentage !== undefined && q.percentage > 80,
			);
			overview.quota_summary = {
				used: usage.length,
				total: usage.length,
				critical_quotas: critical,
			};
		}

		// Current usage
		if (currentUsage.status === "fulfilled") {
			overview.current_usage = currentUsage.value;
		}

		// Billing status
		if (paymentMethods.status === "fulfilled") {
			const methods = paymentMethods.value;
			const primary = methods.find((m) => m.is_primary);
			overview.billing_status = {
				payment_method_status: primary?.status ?? "NO_PAYMENT_METHOD",
			};
		}

		return overview;
	}
}

/**
 * Cached client instance
 */
let cachedClient: SubscriptionClient | null = null;

/**
 * Get or create a Subscription client instance
 *
 * Uses lazy initialization pattern to avoid creating client until needed
 *
 * @param apiClient - The API client from the session
 * @returns The Subscription client instance
 */
export function getSubscriptionClient(
	apiClient: APIClient,
): SubscriptionClient {
	if (!cachedClient) {
		cachedClient = new SubscriptionClient(apiClient);
	}
	return cachedClient;
}

/**
 * Reset the cached client (for testing or session changes)
 */
export function resetSubscriptionClient(): void {
	cachedClient = null;
}
