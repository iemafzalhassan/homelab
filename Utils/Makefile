.PHONY: get-credentials update-ip

get-credentials:
	@echo "Fetching AKS credentials..."
	az aks get-credentials --resource-group homelab-rg --name homelab-aks --overwrite-existing

update-ip:
	@echo "Fetching current public IPv4 address..."
	@CURRENT_IP=$$(curl -4 -s ifconfig.me) ; \
	echo "Current IP: $$CURRENT_IP" ; \
	echo "Updating AKS authorized IP ranges..." ; \
	az aks update --resource-group homelab-rg --name homelab-aks --api-server-authorized-ip-ranges "$$CURRENT_IP/32" ; \
	echo "Updating terraform.tfvars with new IP..." ; \
	sed -i '' "s|admin_ip_range.*|admin_ip_range       = \"$$CURRENT_IP/32\"|g" ../terraform.tfvars ; \
	echo "Done."
