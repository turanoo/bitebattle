// AppTile.swift
// Minimalistic, consistent tile/card for BiteBattle
import SwiftUI


struct AppTile<Content: View>: View {
    let isSelected: Bool
    let action: (() -> Void)?
    let content: () -> Content

    init(isSelected: Bool = false, action: (() -> Void)? = nil, @ViewBuilder content: @escaping () -> Content) {
        self.isSelected = isSelected
        self.action = action
        self.content = content
    }

    var body: some View {
        Group {
            if let action = action {
                Button(action: action) {
                    tileContent
                }
                .buttonStyle(PlainButtonStyle())
            } else {
                tileContent
            }
        }
        .background(isSelected ? AppColors.tileSelected : AppColors.tileBackground)
        .cornerRadius(14)
        .overlay(
            RoundedRectangle(cornerRadius: 14)
                .stroke(isSelected ? AppColors.accent : Color.clear, lineWidth: isSelected ? 2 : 0)
        )
        .shadow(color: Color.black.opacity(0.04), radius: 2, x: 0, y: 1)
        .animation(.easeInOut(duration: 0.18), value: isSelected)
    }

    private var tileContent: some View {
        content()
            .padding()
    }
}
