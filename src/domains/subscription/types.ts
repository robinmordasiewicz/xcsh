/**
 * Subscription Domain Types
 *
 * Type definitions for F5 Distributed Cloud subscription, billing, and quota APIs
 */

/**
 * Common metadata for F5 XC objects
 */
export interface ObjectMeta {
	name: string;
	namespace?: string;
	uid?: string;
	description?: string;
	labels?: Record<string, string>;
	annotations?: Record<string, string>;
	creation_timestamp?: string;
	modification_timestamp?: string;
}

/**
 * Usage Plan - Subscription tier information
 */
export interface UsagePlan {
	metadata?: ObjectMeta;
	name: string;
	display_name?: string;
	description?: string;
	plan_type?: string;
	tier?: string;
	features?: string[];
	limits?: Record<string, number>;
	status?: string;
}

/**
 * Usage Plan List Response
 */
export interface UsagePlanListResponse {
	items?: UsagePlan[];
}

/**
 * Addon Service - Available addon features
 */
export interface AddonService {
	metadata?: ObjectMeta;
	name: string;
	display_name?: string;
	description?: string;
	category?: string;
	pricing?: AddonPricing;
	features?: string[];
	status?: "AVAILABLE" | "COMING_SOON" | "DEPRECATED";
	requires?: string[];
	access_type?: "ENABLED" | "DISABLED" | "TRIAL" | "SUBSCRIBED";
}

/**
 * Addon Pricing Information
 */
export interface AddonPricing {
	model?: "USAGE_BASED" | "FLAT_RATE" | "TIERED";
	base_price?: number;
	currency?: string;
	billing_period?: "MONTHLY" | "YEARLY";
	unit?: string;
	unit_price?: number;
}

/**
 * Addon Service List Response
 */
export interface AddonServiceListResponse {
	items?: AddonService[];
}

/**
 * Addon Subscription - User's addon subscriptions
 */
export interface AddonSubscription {
	metadata?: ObjectMeta;
	name: string;
	addon_service?: string;
	status?: "ACTIVE" | "PENDING" | "CANCELLED" | "EXPIRED";
	start_date?: string;
	end_date?: string;
	auto_renew?: boolean;
	usage?: AddonUsage;
}

/**
 * Addon Usage Information
 */
export interface AddonUsage {
	current_usage?: number;
	limit?: number;
	unit?: string;
	period_start?: string;
	period_end?: string;
}

/**
 * Addon Subscription List Response
 */
export interface AddonSubscriptionListResponse {
	items?: AddonSubscription[];
}

/**
 * Addon Activation Status
 */
export interface AddonActivationStatus {
	addon_name: string;
	status: "ACTIVE" | "INACTIVE" | "PENDING" | "FAILED";
	activated_at?: string;
	error_message?: string;
}

/**
 * Quota Limit - Single quota limit definition
 */
export interface QuotaLimit {
	name: string;
	display_name?: string;
	description?: string;
	limit: number;
	unit?: string;
	category?: string;
	scope?: "TENANT" | "NAMESPACE";
}

/**
 * Quota Limits Response
 */
export interface QuotaLimitsResponse {
	limits?: QuotaLimit[];
	tenant_name?: string;
	plan_type?: string;
}

/**
 * Quota Usage - Current usage against limits
 */
export interface QuotaUsage {
	name: string;
	display_name?: string;
	current: number;
	limit: number;
	unit?: string;
	percentage?: number;
	category?: string;
}

/**
 * Quota Usage Response
 */
export interface QuotaUsageResponse {
	usage?: QuotaUsage[];
	tenant_name?: string;
	as_of?: string;
}

/**
 * Current Usage - Active billing period usage
 */
export interface CurrentUsage {
	billing_period_start?: string;
	billing_period_end?: string;
	total_cost?: number;
	currency?: string;
	usage_items?: UsageItem[];
	projected_cost?: number;
}

/**
 * Usage Item - Individual usage line item
 */
export interface UsageItem {
	name: string;
	display_name?: string;
	category?: string;
	quantity: number;
	unit?: string;
	unit_price?: number;
	total_cost?: number;
}

/**
 * Hourly Usage Details
 */
export interface HourlyUsageDetails {
	date: string;
	hours?: HourlyUsageEntry[];
}

/**
 * Hourly Usage Entry
 */
export interface HourlyUsageEntry {
	hour: number;
	usage_items?: UsageItem[];
	total_cost?: number;
}

/**
 * Monthly Usage Summary
 */
export interface MonthlyUsage {
	month: string;
	year: number;
	total_cost?: number;
	currency?: string;
	usage_by_category?: Record<string, number>;
	usage_items?: UsageItem[];
	invoice_id?: string;
	invoice_status?: string;
}

/**
 * Monthly Usage List Response
 */
export interface MonthlyUsageListResponse {
	items?: MonthlyUsage[];
}

/**
 * Payment Method
 */
export interface PaymentMethod {
	metadata?: ObjectMeta;
	name: string;
	type: "CREDIT_CARD" | "BANK_ACCOUNT" | "INVOICE" | "OTHER";
	is_primary?: boolean;
	is_secondary?: boolean;
	last_four?: string;
	expiry_month?: number;
	expiry_year?: number;
	card_brand?: string;
	billing_address?: BillingAddress;
	status?: "ACTIVE" | "EXPIRED" | "INVALID";
}

/**
 * Billing Address
 */
export interface BillingAddress {
	street?: string;
	city?: string;
	state?: string;
	postal_code?: string;
	country?: string;
}

/**
 * Payment Method List Response
 */
export interface PaymentMethodListResponse {
	items?: PaymentMethod[];
}

/**
 * Invoice
 */
export interface Invoice {
	metadata?: ObjectMeta;
	invoice_id: string;
	invoice_number?: string;
	billing_period_start?: string;
	billing_period_end?: string;
	issue_date?: string;
	due_date?: string;
	total_amount?: number;
	currency?: string;
	status?: "DRAFT" | "OPEN" | "PAID" | "OVERDUE" | "VOID";
	line_items?: InvoiceLineItem[];
	payment_method?: string;
	pdf_url?: string;
}

/**
 * Invoice Line Item
 */
export interface InvoiceLineItem {
	description: string;
	quantity?: number;
	unit_price?: number;
	amount?: number;
	category?: string;
}

/**
 * Invoice List Response
 */
export interface InvoiceListResponse {
	items?: Invoice[];
}

/**
 * Plan Transition Request
 */
export interface PlanTransitionRequest {
	target_plan: string;
	effective_date?: string;
	reason?: string;
}

/**
 * Plan Transition Response
 */
export interface PlanTransitionResponse {
	transition_id?: string;
	status?: "PENDING" | "APPROVED" | "REJECTED" | "COMPLETED";
	current_plan?: string;
	target_plan?: string;
	effective_date?: string;
	message?: string;
}

/**
 * Subscribe to Addon Request
 */
export interface AddonSubscribeRequest {
	addon_service: string;
	billing_cycle?: "MONTHLY" | "YEARLY";
	auto_renew?: boolean;
}

/**
 * Unsubscribe from Addon Request
 */
export interface AddonUnsubscribeRequest {
	addon_service: string;
	reason?: string;
	immediate?: boolean;
}

/**
 * Subscription Overview - Combined summary for default command
 */
export interface SubscriptionOverview {
	plan?: UsagePlan;
	addon_summary?: {
		active: number;
		available: number;
		total: number;
	};
	quota_summary?: {
		used: number;
		total: number;
		critical_quotas?: QuotaUsage[];
	};
	current_usage?: CurrentUsage;
	billing_status?: {
		next_invoice_date?: string;
		outstanding_balance?: number;
		payment_method_status?: string;
	};
}

/**
 * Report Summary - Comprehensive subscription report
 */
export interface SubscriptionReport {
	generated_at: string;
	tenant_name?: string;
	plan: UsagePlan;
	addons: {
		subscribed: AddonSubscription[];
		available: AddonService[];
	};
	quotas: {
		limits: QuotaLimit[];
		usage: QuotaUsage[];
		utilization_percentage: number;
	};
	usage: {
		current: CurrentUsage;
		monthly_trend?: MonthlyUsage[];
	};
	billing: {
		payment_methods: PaymentMethod[];
		recent_invoices: Invoice[];
		total_ytd?: number;
	};
}
