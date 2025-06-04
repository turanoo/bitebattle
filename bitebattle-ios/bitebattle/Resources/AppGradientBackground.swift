// AppGradientBackground.swift
// Reusable gradient background for BiteBattle
import SwiftUI

public struct AppGradientBackground<Content: View>: View {
    let content: Content
    public init(@ViewBuilder content: () -> Content) {
        self.content = content()
    }
    public var body: some View {
        ZStack {
            LinearGradient(
                gradient: Gradient(colors: [Color("#fff7e3"), Color("#ffa43d").opacity(0.7)]),
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            )
            .ignoresSafeArea()
            content
        }
    }
}
