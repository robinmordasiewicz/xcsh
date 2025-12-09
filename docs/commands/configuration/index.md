---
title: "vesctl configuration"
description: "Manage F5 XC configuration objects using CRUD operations."
keywords:
  - F5 Distributed Cloud
  - configuration
  - vesctl
  - F5 XC
command: "vesctl configuration"
command_group: "configuration"
aliases:
  - "cfg"
  - "c"
---

# vesctl configuration

> Manage F5 XC configuration objects using CRUD operations.

## Synopsis

```bash
vesctl configuration <command> [flags]
```

## Aliases

This command can also be invoked as:

- `vesctl cfg`
- `vesctl c`

## Available Commands

| Command | Description |
|---------|-------------|
| [address_allocator](address_allocator.md) | Manage address allocator resources |
| [advertise_policy](advertise_policy.md) | Manage advertise policy resources |
| [ai_assistant](ai_assistant.md) | Manage ai assistant resources |
| [ai_data_bfdp](ai_data_bfdp.md) | Manage ai data bfdp resources |
| [ai_data_bfdp_subscription](ai_data_bfdp_subscription.md) | Manage ai data bfdp subscription resources |
| [alert](alert.md) | Manage alert resources |
| [alert_policy](alert_policy.md) | Manage alert policy resources |
| [alert_receiver](alert_receiver.md) | Manage alert receiver resources |
| [api_credential](api_credential.md) | Manage api credential resources |
| [api_definition](api_definition.md) | Manage api definition resources |
| [api_group](api_group.md) | Manage api group resources |
| [api_group_element](api_group_element.md) | Manage api group element resources |
| [api_sec_api_crawler](api_sec_api_crawler.md) | Manage api sec api crawler resources |
| [api_sec_api_discovery](api_sec_api_discovery.md) | Manage api sec api discovery resources |
| [api_sec_api_testing](api_sec_api_testing.md) | Manage api sec api testing resources |
| [api_sec_code_base_integration](api_sec_code_base_integration.md) | Manage api sec code base integration resources |
| [api_sec_rule_suggestion](api_sec_rule_suggestion.md) | Manage api sec rule suggestion resources |
| [app_api_group](app_api_group.md) | Manage app api group resources |
| [app_firewall](app_firewall.md) | Manage app firewall resources |
| [app_security](app_security.md) | Manage app security resources |
| [app_setting](app_setting.md) | Manage app setting resources |
| [app_type](app_type.md) | Manage app type resources |
| [authentication](authentication.md) | Manage authentication resources |
| [aws_tgw_site](aws_tgw_site.md) | Manage aws tgw site resources |
| [aws_vpc_site](aws_vpc_site.md) | Manage aws vpc site resources |
| [azure_vnet_site](azure_vnet_site.md) | Manage azure vnet site resources |
| [bgp](bgp.md) | Manage bgp resources |
| [bgp_asn_set](bgp_asn_set.md) | Manage bgp asn set resources |
| [bgp_routing_policy](bgp_routing_policy.md) | Manage bgp routing policy resources |
| [bigcne_data_group](bigcne_data_group.md) | Manage bigcne data group resources |
| [bigcne_irule](bigcne_irule.md) | Manage bigcne irule resources |
| [bigip_apm](bigip_apm.md) | Manage bigip apm resources |
| [bigip_irule](bigip_irule.md) | Manage bigip irule resources |
| [bigip_virtual_server](bigip_virtual_server.md) | Manage bigip virtual server resources |
| [billing_payment_method](billing_payment_method.md) | Manage billing payment method resources |
| [billing_plan_transition](billing_plan_transition.md) | Manage billing plan transition resources |
| [bot_defense_app_infrastructure](bot_defense_app_infrastructure.md) | Manage bot defense app infrastructure resources |
| [cdn_cache_rule](cdn_cache_rule.md) | Manage cdn cache rule resources |
| [cdn_loadbalancer](cdn_loadbalancer.md) | Manage cdn loadbalancer resources |
| [certificate](certificate.md) | Manage certificate resources |
| [certificate_chain](certificate_chain.md) | Manage certificate chain resources |
| [certified_hardware](certified_hardware.md) | Manage certified hardware resources |
| [cloud_connect](cloud_connect.md) | Manage cloud connect resources |
| [cloud_credentials](cloud_credentials.md) | Manage cloud credentials resources |
| [cloud_elastic_ip](cloud_elastic_ip.md) | Manage cloud elastic ip resources |
| [cloud_link](cloud_link.md) | Manage cloud link resources |
| [cloud_region](cloud_region.md) | Manage cloud region resources |
| [cluster](cluster.md) | Manage cluster resources |
| [cminstance](cminstance.md) | Manage cminstance resources |
| [contact](contact.md) | Manage contact resources |
| [container_registry](container_registry.md) | Manage container registry resources |
| [crl](crl.md) | Manage crl resources |
| [customer_support](customer_support.md) | Manage customer support resources |
| [data_privacy_geo_config](data_privacy_geo_config.md) | Manage data privacy geo config resources |
| [data_privacy_lma_region](data_privacy_lma_region.md) | Manage data privacy lma region resources |
| [data_type](data_type.md) | Manage data type resources |
| [dc_cluster_group](dc_cluster_group.md) | Manage dc cluster group resources |
| [discovered_service](discovered_service.md) | Manage discovered service resources |
| [discovery](discovery.md) | Manage discovery resources |
| [dns_compliance_checks](dns_compliance_checks.md) | Manage dns compliance checks resources |
| [dns_domain](dns_domain.md) | Manage dns domain resources |
| [dns_lb_health_check](dns_lb_health_check.md) | Manage dns lb health check resources |
| [dns_lb_pool](dns_lb_pool.md) | Manage dns lb pool resources |
| [dns_load_balancer](dns_load_balancer.md) | Manage dns load balancer resources |
| [dns_zone](dns_zone.md) | Manage dns zone resources |
| [dns_zone_rrset](dns_zone_rrset.md) | Manage dns zone rrset resources |
| [dns_zone_subscription](dns_zone_subscription.md) | Manage dns zone subscription resources |
| [endpoint](endpoint.md) | Manage endpoint resources |
| [enhanced_firewall_policy](enhanced_firewall_policy.md) | Manage enhanced firewall policy resources |
| [external_connector](external_connector.md) | Manage external connector resources |
| [fast_acl](fast_acl.md) | Manage fast acl resources |
| [fast_acl_rule](fast_acl_rule.md) | Manage fast acl rule resources |
| [filter_set](filter_set.md) | Manage filter set resources |
| [fleet](fleet.md) | Manage fleet resources |
| [flow](flow.md) | Manage flow resources |
| [flow_anomaly](flow_anomaly.md) | Manage flow anomaly resources |
| [forward_proxy_policy](forward_proxy_policy.md) | Manage forward proxy policy resources |
| [forwarding_class](forwarding_class.md) | Manage forwarding class resources |
| [gcp_vpc_site](gcp_vpc_site.md) | Manage gcp vpc site resources |
| [geo_location_set](geo_location_set.md) | Manage geo location set resources |
| [gia](gia.md) | Manage gia resources |
| [global_log_receiver](global_log_receiver.md) | Manage global log receiver resources |
| [graph_connectivity](graph_connectivity.md) | Manage graph connectivity resources |
| [graph_l3l4](graph_l3l4.md) | Manage graph l3l4 resources |
| [graph_service](graph_service.md) | Manage graph service resources |
| [graph_site](graph_site.md) | Manage graph site resources |
| [healthcheck](healthcheck.md) | Manage healthcheck resources |
| [http_loadbalancer](http_loadbalancer.md) | Manage http loadbalancer resources |
| [ike1](ike1.md) | Manage ike1 resources |
| [ike2](ike2.md) | Manage ike2 resources |
| [ike_phase1_profile](ike_phase1_profile.md) | Manage ike phase1 profile resources |
| [ike_phase2_profile](ike_phase2_profile.md) | Manage ike phase2 profile resources |
| [implicit_label](implicit_label.md) | Manage implicit label resources |
| [infraprotect](infraprotect.md) | Manage infraprotect resources |
| [infraprotect_asn](infraprotect_asn.md) | Manage infraprotect asn resources |
| [infraprotect_asn_prefix](infraprotect_asn_prefix.md) | Manage infraprotect asn prefix resources |
| [infraprotect_deny_list_rule](infraprotect_deny_list_rule.md) | Manage infraprotect deny list rule resources |
| [infraprotect_firewall_rule](infraprotect_firewall_rule.md) | Manage infraprotect firewall rule resources |
| [infraprotect_firewall_rule_group](infraprotect_firewall_rule_group.md) | Manage infraprotect firewall rule group resources |
| [infraprotect_firewall_ruleset](infraprotect_firewall_ruleset.md) | Manage infraprotect firewall ruleset resources |
| [infraprotect_information](infraprotect_information.md) | Manage infraprotect information resources |
| [infraprotect_internet_prefix_advertisement](infraprotect_internet_prefix_advertisement.md) | Manage infraprotect internet prefix advertisement resources |
| [infraprotect_tunnel](infraprotect_tunnel.md) | Manage infraprotect tunnel resources |
| [ip_prefix_set](ip_prefix_set.md) | Manage ip prefix set resources |
| [k8s_cluster](k8s_cluster.md) | Manage k8s cluster resources |
| [k8s_cluster_role](k8s_cluster_role.md) | Manage k8s cluster role resources |
| [k8s_cluster_role_binding](k8s_cluster_role_binding.md) | Manage k8s cluster role binding resources |
| [k8s_pod_security_admission](k8s_pod_security_admission.md) | Manage k8s pod security admission resources |
| [k8s_pod_security_policy](k8s_pod_security_policy.md) | Manage k8s pod security policy resources |
| [known_label](known_label.md) | Manage known label resources |
| [known_label_key](known_label_key.md) | Manage known label key resources |
| [log](log.md) | Manage log resources |
| [log_receiver](log_receiver.md) | Manage log receiver resources |
| [maintenance_status](maintenance_status.md) | Manage maintenance status resources |
| [malicious_user_mitigation](malicious_user_mitigation.md) | Manage malicious user mitigation resources |
| [malware_protection_subscription](malware_protection_subscription.md) | Manage malware protection subscription resources |
| [marketplace_aws_account](marketplace_aws_account.md) | Manage marketplace aws account resources |
| [marketplace_xc_saas](marketplace_xc_saas.md) | Manage marketplace xc saas resources |
| [module_management](module_management.md) | Manage module management resources |
| [namespace](namespace.md) | Manage namespace resources |
| [namespace_role](namespace_role.md) | Manage namespace role resources |
| [nat_policy](nat_policy.md) | Manage nat policy resources |
| [network_connector](network_connector.md) | Manage network connector resources |
| [network_firewall](network_firewall.md) | Manage network firewall resources |
| [network_interface](network_interface.md) | Manage network interface resources |
| [network_policy](network_policy.md) | Manage network policy resources |
| [network_policy_rule](network_policy_rule.md) | Manage network policy rule resources |
| [network_policy_set](network_policy_set.md) | Manage network policy set resources |
| [network_policy_view](network_policy_view.md) | Manage network policy view resources |
| [nfv_service](nfv_service.md) | Manage nfv service resources |
| [nginx_one_nginx_csg](nginx_one_nginx_csg.md) | Manage nginx one nginx csg resources |
| [nginx_one_nginx_instance](nginx_one_nginx_instance.md) | Manage nginx one nginx instance resources |
| [nginx_one_nginx_server](nginx_one_nginx_server.md) | Manage nginx one nginx server resources |
| [nginx_one_nginx_service_discovery](nginx_one_nginx_service_discovery.md) | Manage nginx one nginx service discovery resources |
| [nginx_one_subscription](nginx_one_subscription.md) | Manage nginx one subscription resources |
| [observability_subscription](observability_subscription.md) | Manage observability subscription resources |
| [oidc_provider](oidc_provider.md) | Manage oidc provider resources |
| [operate_bgp](operate_bgp.md) | Manage operate bgp resources |
| [operate_crl](operate_crl.md) | Manage operate crl resources |
| [operate_debug](operate_debug.md) | Manage operate debug resources |
| [operate_dhcp](operate_dhcp.md) | Manage operate dhcp resources |
| [operate_flow](operate_flow.md) | Manage operate flow resources |
| [operate_lte](operate_lte.md) | Manage operate lte resources |
| [operate_ping](operate_ping.md) | Manage operate ping resources |
| [operate_route](operate_route.md) | Manage operate route resources |
| [operate_tcpdump](operate_tcpdump.md) | Manage operate tcpdump resources |
| [operate_traceroute](operate_traceroute.md) | Manage operate traceroute resources |
| [operate_usb](operate_usb.md) | Manage operate usb resources |
| [operate_wifi](operate_wifi.md) | Manage operate wifi resources |
| [origin_pool](origin_pool.md) | Manage origin pool resources |
| [pbac_addon_service](pbac_addon_service.md) | Manage pbac addon service resources |
| [pbac_addon_subscription](pbac_addon_subscription.md) | Manage pbac addon subscription resources |
| [pbac_catalog](pbac_catalog.md) | Manage pbac catalog resources |
| [pbac_navigation_tile](pbac_navigation_tile.md) | Manage pbac navigation tile resources |
| [pbac_plan](pbac_plan.md) | Manage pbac plan resources |
| [policer](policer.md) | Manage policer resources |
| [policy_based_routing](policy_based_routing.md) | Manage policy based routing resources |
| [protocol_inspection](protocol_inspection.md) | Manage protocol inspection resources |
| [protocol_policer](protocol_policer.md) | Manage protocol policer resources |
| [proxy](proxy.md) | Manage proxy resources |
| [public_ip](public_ip.md) | Manage public ip resources |
| [quota](quota.md) | Manage quota resources |
| [rate_limiter](rate_limiter.md) | Manage rate limiter resources |
| [rate_limiter_policy](rate_limiter_policy.md) | Manage rate limiter policy resources |
| [rbac_policy](rbac_policy.md) | Manage rbac policy resources |
| [registration](registration.md) | Manage registration resources |
| [report](report.md) | Manage report resources |
| [report_config](report_config.md) | Manage report config resources |
| [role](role.md) | Manage role resources |
| [route](route.md) | Manage route resources |
| [scim](scim.md) | Manage scim resources |
| [secret_management](secret_management.md) | Manage secret management resources |
| [secret_management_access](secret_management_access.md) | Manage secret management access resources |
| [secret_policy](secret_policy.md) | Manage secret policy resources |
| [secret_policy_rule](secret_policy_rule.md) | Manage secret policy rule resources |
| [securemesh_site](securemesh_site.md) | Manage securemesh site resources |
| [securemesh_site_v2](securemesh_site_v2.md) | Manage securemesh site v2 resources |
| [segment](segment.md) | Manage segment resources |
| [segment_connection](segment_connection.md) | Manage segment connection resources |
| [sensitive_data_policy](sensitive_data_policy.md) | Manage sensitive data policy resources |
| [service_policy](service_policy.md) | Manage service policy resources |
| [service_policy_rule](service_policy_rule.md) | Manage service policy rule resources |
| [service_policy_set](service_policy_set.md) | Manage service policy set resources |
| [shape_bot_defense_bot_allowlist_policy](shape_bot_defense_bot_allowlist_policy.md) | Manage shape bot defense bot allowlist policy resources |
| [shape_bot_defense_bot_endpoint_policy](shape_bot_defense_bot_endpoint_policy.md) | Manage shape bot defense bot endpoint policy resources |
| [shape_bot_defense_bot_infrastructure](shape_bot_defense_bot_infrastructure.md) | Manage shape bot defense bot infrastructure resources |
| [shape_bot_defense_bot_network_policy](shape_bot_defense_bot_network_policy.md) | Manage shape bot defense bot network policy resources |
| [shape_bot_defense_instance](shape_bot_defense_instance.md) | Manage shape bot defense instance resources |
| [shape_bot_defense_mobile_base_config](shape_bot_defense_mobile_base_config.md) | Manage shape bot defense mobile base config resources |
| [shape_bot_defense_mobile_sdk](shape_bot_defense_mobile_sdk.md) | Manage shape bot defense mobile sdk resources |
| [shape_bot_defense_protected_application](shape_bot_defense_protected_application.md) | Manage shape bot defense protected application resources |
| [shape_bot_defense_reporting](shape_bot_defense_reporting.md) | Manage shape bot defense reporting resources |
| [shape_bot_defense_subscription](shape_bot_defense_subscription.md) | Manage shape bot defense subscription resources |
| [shape_bot_detection_rule](shape_bot_detection_rule.md) | Manage shape bot detection rule resources |
| [shape_bot_detection_update](shape_bot_detection_update.md) | Manage shape bot detection update resources |
| [shape_brmalerts_alert_gen_policy](shape_brmalerts_alert_gen_policy.md) | Manage shape brmalerts alert gen policy resources |
| [shape_brmalerts_alert_template](shape_brmalerts_alert_template.md) | Manage shape brmalerts alert template resources |
| [shape_client_side_defense](shape_client_side_defense.md) | Manage shape client side defense resources |
| [shape_client_side_defense_allowed_domain](shape_client_side_defense_allowed_domain.md) | Manage shape client side defense allowed domain resources |
| [shape_client_side_defense_mitigated_domain](shape_client_side_defense_mitigated_domain.md) | Manage shape client side defense mitigated domain resources |
| [shape_client_side_defense_protected_domain](shape_client_side_defense_protected_domain.md) | Manage shape client side defense protected domain resources |
| [shape_client_side_defense_subscription](shape_client_side_defense_subscription.md) | Manage shape client side defense subscription resources |
| [shape_data_delivery](shape_data_delivery.md) | Manage shape data delivery resources |
| [shape_data_delivery_receiver](shape_data_delivery_receiver.md) | Manage shape data delivery receiver resources |
| [shape_data_delivery_subscription](shape_data_delivery_subscription.md) | Manage shape data delivery subscription resources |
| [shape_device_id](shape_device_id.md) | Manage shape device id resources |
| [shape_mobile_app_shield_subscription](shape_mobile_app_shield_subscription.md) | Manage shape mobile app shield subscription resources |
| [shape_mobile_integrator_subscription](shape_mobile_integrator_subscription.md) | Manage shape mobile integrator subscription resources |
| [shape_recognize](shape_recognize.md) | Manage shape recognize resources |
| [shape_safe](shape_safe.md) | Manage shape safe resources |
| [shape_safeap](shape_safeap.md) | Manage shape safeap resources |
| [signup](signup.md) | Manage signup resources |
| [site](site.md) | Manage site resources |
| [site_mesh_group](site_mesh_group.md) | Manage site mesh group resources |
| [srv6_network_slice](srv6_network_slice.md) | Manage srv6 network slice resources |
| [status_at_site](status_at_site.md) | Manage status at site resources |
| [stored_object](stored_object.md) | Manage stored object resources |
| [subnet](subnet.md) | Manage subnet resources |
| [subscription](subscription.md) | Manage subscription resources |
| [synthetic_monitor](synthetic_monitor.md) | Manage synthetic monitor resources |
| [synthetic_monitor_dns](synthetic_monitor_dns.md) | Manage synthetic monitor dns resources |
| [synthetic_monitor_http](synthetic_monitor_http.md) | Manage synthetic monitor http resources |
| [tcp_loadbalancer](tcp_loadbalancer.md) | Manage tcp loadbalancer resources |
| [tenant](tenant.md) | Manage tenant resources |
| [tenant_configuration](tenant_configuration.md) | Manage tenant configuration resources |
| [tenant_management](tenant_management.md) | Manage tenant management resources |
| [tenant_management_allowed_tenant](tenant_management_allowed_tenant.md) | Manage tenant management allowed tenant resources |
| [tenant_management_child_tenant](tenant_management_child_tenant.md) | Manage tenant management child tenant resources |
| [tenant_management_child_tenant_manager](tenant_management_child_tenant_manager.md) | Manage tenant management child tenant manager resources |
| [tenant_management_managed_tenant](tenant_management_managed_tenant.md) | Manage tenant management managed tenant resources |
| [tenant_management_tenant_profile](tenant_management_tenant_profile.md) | Manage tenant management tenant profile resources |
| [third_party_application](third_party_application.md) | Manage third party application resources |
| [ticket_management_ticket_tracking_system](ticket_management_ticket_tracking_system.md) | Manage ticket management ticket tracking system resources |
| [token](token.md) | Manage token resources |
| [topology](topology.md) | Manage topology resources |
| [tpm_api_key](tpm_api_key.md) | Manage tpm api key resources |
| [tpm_category](tpm_category.md) | Manage tpm category resources |
| [tpm_manager](tpm_manager.md) | Manage tpm manager resources |
| [tpm_provision](tpm_provision.md) | Manage tpm provision resources |
| [trusted_ca_list](trusted_ca_list.md) | Manage trusted ca list resources |
| [tunnel](tunnel.md) | Manage tunnel resources |
| [udp_loadbalancer](udp_loadbalancer.md) | Manage udp loadbalancer resources |
| [ui_static_component](ui_static_component.md) | Manage ui static component resources |
| [upgrade_status](upgrade_status.md) | Manage upgrade status resources |
| [usage](usage.md) | Manage usage resources |
| [usage_invoice](usage_invoice.md) | Manage usage invoice resources |
| [usage_plan](usage_plan.md) | Manage usage plan resources |
| [usb_policy](usb_policy.md) | Manage usb policy resources |
| [user](user.md) | Manage user resources |
| [user_group](user_group.md) | Manage user group resources |
| [user_identification](user_identification.md) | Manage user identification resources |
| [user_setting](user_setting.md) | Manage user setting resources |
| [views_terraform_parameters](views_terraform_parameters.md) | Manage views terraform parameters resources |
| [views_view_internal](views_view_internal.md) | Manage views view internal resources |
| [virtual_appliance](virtual_appliance.md) | Manage virtual appliance resources |
| [virtual_host](virtual_host.md) | Manage virtual host resources |
| [virtual_k8s](virtual_k8s.md) | Manage virtual k8s resources |
| [virtual_network](virtual_network.md) | Manage virtual network resources |
| [virtual_site](virtual_site.md) | Manage virtual site resources |
| [voltshare](voltshare.md) | Manage voltshare resources |
| [voltshare_admin_policy](voltshare_admin_policy.md) | Manage voltshare admin policy resources |
| [voltstack_site](voltstack_site.md) | Manage voltstack site resources |
| [waf](waf.md) | Manage waf resources |
| [waf_exclusion_policy](waf_exclusion_policy.md) | Manage waf exclusion policy resources |
| [waf_signatures_changelog](waf_signatures_changelog.md) | Manage waf signatures changelog resources |
| [was_user_token](was_user_token.md) | Manage was user token resources |
| [workload](workload.md) | Manage workload resources |
| [workload_flavor](workload_flavor.md) | Manage workload flavor resources |

## Examples

```bash
vesctl configuration create virtual_host
```

## See Also

- [Command Reference](../index.md)
- [vesctl configuration address_allocator](address_allocator.md)
- [vesctl configuration advertise_policy](advertise_policy.md)
- [vesctl configuration ai_assistant](ai_assistant.md)
- [vesctl configuration ai_data_bfdp](ai_data_bfdp.md)
- [vesctl configuration ai_data_bfdp_subscription](ai_data_bfdp_subscription.md)
