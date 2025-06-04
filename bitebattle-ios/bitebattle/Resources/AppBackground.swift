// AppBackground.swift
import SwiftUI

public struct AppBackground<Content: View>: View {
    let content: Content
    public init(@ViewBuilder content: () -> Content) {
        self.content = content()
    }
    public var body: some View {
        ZStack {
            AppColors.background.ignoresSafeArea()
            content
        }
    }
}
