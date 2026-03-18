import 'package:flutter/material.dart';

class AppTheme {
  static const _accent = Color(0xFF00C9A7);
  static const _accentDark = Color(0xFF00A88A);
  static const _backgroundLight = Color(0xFFF8FAFC);
  static const _backgroundDark = Color(0xFF0F172A);
  static const _surfaceDark = Color(0xFF1E293B);
  static const _vitonGreen = Color(0xFF00E676);
  static const _vitonBlack = Color(0xFF000000);

  static ThemeData get light {
    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.light,
      colorScheme: ColorScheme.light(
        primary: _accent,
        surface: _backgroundLight,
        onPrimary: Colors.white,
        onSurface: const Color(0xFF1E293B),
      ),
      scaffoldBackgroundColor: _backgroundLight,
      appBarTheme: const AppBarTheme(
        elevation: 0,
        centerTitle: true,
        backgroundColor: _backgroundLight,
        foregroundColor: Color(0xFF1E293B),
      ),
      cardTheme: CardThemeData(
        elevation: 0,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        color: Colors.white,
      ),
    );
  }

  static ThemeData get dark {
    return ThemeData(
      useMaterial3: true,
      brightness: Brightness.dark,
      colorScheme: ColorScheme.dark(
        primary: _vitonGreen,
        surface: _vitonBlack,
        onPrimary: Colors.black,
        onSurface: Colors.white,
      ),
      scaffoldBackgroundColor: _vitonBlack,
      appBarTheme: const AppBarTheme(
        elevation: 0,
        centerTitle: true,
        backgroundColor: _vitonBlack,
        foregroundColor: Colors.white,
      ),
      cardTheme: CardThemeData(
        elevation: 0,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        color: _surfaceDark,
      ),
    );
  }

  static const Color vitonGreen = _vitonGreen;
  static const Color vitonBlack = _vitonBlack;
}
