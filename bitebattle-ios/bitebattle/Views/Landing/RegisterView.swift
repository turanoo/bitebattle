import SwiftUI

struct RegisterView: View {
    @State private var email = ""
    @State private var password = ""
    @State private var name = ""
    @State private var registrationStatus: String?
    @AppStorage("isLoggedIn") var isLoggedIn: Bool = false

    var body: some View {
        ZStack {
            // Background gradient
            LinearGradient(
                gradient: Gradient(colors: [Color.pink.opacity(0.7), Color.orange.opacity(0.7)]),
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            )
            .ignoresSafeArea()

            VStack(spacing: 28) {
                Spacer()
                
                // Logo
                Image(systemName: "fork.knife.circle.fill")
                    .resizable()
                    .scaledToFit()
                    .frame(width: 70, height: 70)
                    .foregroundColor(.white)
                    .shadow(radius: 8)

                Text("Create Account")
                    .font(.largeTitle)
                    .fontWeight(.bold)
                    .foregroundColor(.white)
                    .shadow(radius: 3)

                VStack(spacing: 16) {
                    TextField("Full Name", text: $name)
                        .padding()
                        .background(Color.white.opacity(0.9))
                        .cornerRadius(10)
                        .autocapitalization(.words)
                        .disableAutocorrection(true)
                        .shadow(radius: 1)

                    TextField("Email", text: $email)
                        .padding()
                        .background(Color.white.opacity(0.9))
                        .cornerRadius(10)
                        .keyboardType(.emailAddress)
                        .autocapitalization(.none)
                        .disableAutocorrection(true)
                        .shadow(radius: 1)

                    SecureField("Password", text: $password)
                        .padding()
                        .background(Color.white.opacity(0.9))
                        .cornerRadius(10)
                        .shadow(radius: 1)
                }
                .padding(.horizontal, 24)

                Button(action: registerUser) {
                    Text("Register")
                        .fontWeight(.semibold)
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(Color.white.opacity(0.9))
                        .foregroundColor(.pink)
                        .cornerRadius(12)
                        .shadow(radius: 2)
                }
                .padding(.horizontal, 24)

                if let status = registrationStatus {
                    Text(status)
                        .foregroundColor(.white)
                        .font(.footnote)
                        .padding(.top, 8)
                        .shadow(radius: 1)
                }

                Spacer()
            }
            .padding()
        }
    }

    func registerUser() {
        guard let url = URL(string: "http://localhost:8080/api/register") else { return }

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let payload: [String: String] = [
            "email": email,
            "password": password,
            "name": name
        ]

        guard let httpBody = try? JSONSerialization.data(withJSONObject: payload) else { return }

        request.httpBody = httpBody

        URLSession.shared.dataTask(with: request) { data, response, error in
            DispatchQueue.main.async {
                if let error = error {
                    registrationStatus = "Failed: \(error.localizedDescription)"
                    return
                }

                if let httpResponse = response as? HTTPURLResponse {
                    if httpResponse.statusCode == 201 {
                        // Parse token from response data
                        if let data = data,
                           let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
                           let token = json["token"] as? String {
                            UserDefaults.standard.set(token, forKey: "authToken")
                            isLoggedIn = true // Trigger navigation to HomeView
                            registrationStatus = "Successfully registered!"
                        } else {
                            registrationStatus = "Registration succeeded, but token missing."
                        }
                    } else {
                        registrationStatus = "Failed: Status \(httpResponse.statusCode)"
                    }
                }
            }
        }.resume()
    }
}
