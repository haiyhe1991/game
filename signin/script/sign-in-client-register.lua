----注册方法
register_proto_method("pactum.GatewayRegisterRequest", "logic.GatewayRegisterProc", "OnProccess")
register_proto_method("pactum.LoginRequest",   "logic.SignInProc",  "OnProccess")
register_proto_method("pactum.UnLoginRequest", "logic.SignOutProc", "OnProccess")
