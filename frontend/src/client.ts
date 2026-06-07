import { createConnectTransport } from "@connectrpc/connect-web"
import { createClient } from "@connectrpc/connect"

import { ProviderService } from "@/gen/neoyu/connection/v1/provider_pb"

const transport = createConnectTransport({
  baseUrl: "/api",
})

export const providerClient = createClient(ProviderService, transport)
