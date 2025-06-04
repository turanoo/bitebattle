// AppColors.swift
// Centralized color palette for BiteBattle
import SwiftUI

// Ensure all properties and struct are public for cross-file access
public struct AppColors {
    // Cream/Backgrounds (ordered by your list)
    public static let cream        = Color(hex: "#fff7e3") // Soft off-white / light cream
    public static let creamAlt     = Color(hex: "#fff7e4") // Slight variation of cream
    public static let creamLight   = Color(hex: "#fff8e4") // Light cream tone
    public static let creamSofter  = Color(hex: "#fff8e3") // Slightly different cream
    public static let creamSoftest = Color(hex: "#fff8e5") // Another soft variant
    public static let creamPale    = Color(hex: "#ffffec") // Pale yellow-white
    public static let creamPastel  = Color(hex: "#fffae6") // Light pastel yellow
    public static let creamUltra   = Color(hex: "#fffce8") // Softest cream tone
    public static let creamVar     = Color(hex: "#fff9e5") // Cream variation

    // Brown/Accent
    public static let brown        = Color(hex: "#a0522d") // Warm brown (used in fork)

    // Accent/Primary
    public static let orange   = Color(hex: "#ffa43d") // Warm orange (used in fork/text)
    public static let primary  = orange
    public static let accent   = orange
    public static let secondary = creamPastel // Add this line

    // Surfaces and backgrounds
    public static let surface    = cream
    public static let background = cream

    // Text
    public static let textPrimary    = Color.black
    public static let textOnPrimary  = Color.white
    public static let textSecondary  = Color.gray

    // UI States
    public static let error        = Color.red
    public static let tileBackground = Color.white.opacity(0.12)
    public static let tileSelected   = orange.opacity(0.15)
    public static let disabled       = Color.gray.opacity(0.5)
}

// Helper to use hex colors
public extension Color {
    init(hex: String) {
        let hex = hex.trimmingCharacters(in: CharacterSet.alphanumerics.inverted)
        var int: UInt64 = 0
        Scanner(string: hex).scanHexInt64(&int)
        let a, r, g, b: UInt64
        switch hex.count {
        case 3: // RGB (12-bit)
            (a, r, g, b) = (255, (int >> 8) * 17, (int >> 4 & 0xF) * 17, (int & 0xF) * 17)
        case 6: // RGB (24-bit)
            (a, r, g, b) = (255, int >> 16, int >> 8 & 0xFF, int & 0xFF)
        case 8: // ARGB (32-bit)
            (a, r, g, b) = (int >> 24, int >> 16 & 0xFF, int >> 8 & 0xFF, int & 0xFF)
        default:
            (a, r, g, b) = (255, 255, 255, 255)
        }
        self.init(
            .sRGB,
            red: Double(r) / 255,
            green: Double(g) / 255,
            blue: Double(b) / 255,
            opacity: Double(a) / 255
        )
    }
}
