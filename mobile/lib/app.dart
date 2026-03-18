import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'config/theme.dart';
import 'providers/auth_provider.dart';
import 'screens/home_screen.dart';
import 'screens/login_screen.dart';

class VPNStartupApp extends StatelessWidget {
  const VPNStartupApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Viton VPN',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.light,
      darkTheme: AppTheme.dark,
      themeMode: ThemeMode.dark,
      home: Consumer<AuthProvider>(
        builder: (_, auth, __) {
          if (auth.isLoggedIn) return const HomeScreen();
          return const LoginScreen();
        },
      ),
    );
  }
}
