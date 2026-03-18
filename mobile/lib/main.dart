import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'app.dart';
import 'providers/auth_provider.dart';
import 'providers/vpn_provider.dart';
import 'providers/servers_provider.dart';

void main() {
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => AuthProvider()),
        ChangeNotifierProvider(create: (_) => VPNProvider()),
        ChangeNotifierProvider(create: (_) => ServersProvider()),
      ],
      child: const VPNStartupApp(),
    ),
  );
}
